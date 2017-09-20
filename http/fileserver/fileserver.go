package fileserver

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/prasannavl/go-gluons/http/handlerutils"
	"github.com/prasannavl/goerror"
	"github.com/prasannavl/goerror/httperror"
)

func HttpFileServer(root http.FileSystem) http.Handler {
	return &httpFileHandler{root}
}

type httpFileHandler struct {
	root http.FileSystem
}

func (f *httpFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)
		switch e.Kind {
		case ErrRedirect:
			e.Headers().Write(w)
			w.WriteHeader(http.StatusMovedPermanently)
		default:
			http.Error(w, e.Error(), e.Code())
			//  Alternatively, handle the dir listing.
			// case ErrDirFound
			//	stat := e.Dirstat
			//	finished := HandleDirPrelude(w, r, stat)
			//	if !finished {
			//		// Handle the dir listing here.
			//	}
		}
	}
}

func ServeRequestPath(w http.ResponseWriter, r *http.Request, root http.FileSystem) error {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	return serveFile(w, r, root, path.Clean(upath), true)
}

func ServeFile(w http.ResponseWriter, r *http.Request, name string) error {
	if ContainsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).
		return newErr(http.StatusBadRequest, ErrBadRequest, "invalid URL path", name, nil)
	}
	dir, file := filepath.Split(name)
	return serveFile(w, r, http.Dir(dir), file, false)
}

func HandleDirPrelude(w http.ResponseWriter, r *http.Request, dirstat os.FileInfo) (finished bool) {
	d := dirstat
	if checkIfModifiedSince(r, d.ModTime()) == condFalse {
		writeNotModified(w)
		return true
	}
	w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
	return false
}

func ContainsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

// name is '/'-separated, not filepath.Separator.
func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool) error {
	const indexPage = "/index.html"

	f, err := fs.Open(name)
	if err != nil {
		return newErr(http.StatusNotFound, ErrFsOpen, "", name, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return newErr(http.StatusInternalServerError, ErrFsStat, "", name, err)
	}

	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, indexPage) {
		return newErrRedirect(r, "./", name)
	}

	if redirect {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path
		if stat.IsDir() {
			// redirect if the directory name doesn't end in a slash
			if url[len(url)-1] != '/' {
				return newErrRedirect(r, path.Base(url)+"/", name)
			}
		} else {
			if url[len(url)-1] == '/' {
				return newErrRedirect(r, "../"+path.Base(url), name)
			}
		}
	}

	// use contents of index.html for directory, if present
	if stat.IsDir() {
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := fs.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				name = index
				stat = dd
				f = ff
			}
		}
	}

	// Still a directory? (we didn't find an index.html file)
	if stat.IsDir() {
		return newErrDirFound(name, stat)
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	return nil
}

// condResult is the result of an HTTP request precondition check.
// See https://tools.ietf.org/html/rfc7232 section 3.
type condResult int

const (
	condNone  condResult = iota
	condTrue
	condFalse
)

func checkIfModifiedSince(r *http.Request, modtime time.Time) condResult {
	if r.Method != "GET" && r.Method != "HEAD" {
		return condNone
	}
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" || isZeroTime(modtime) {
		return condNone
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return condNone
	}
	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if modtime.Before(t.Add(1 * time.Second)) {
		return condFalse
	}
	return condTrue
}

func writeNotModified(w http.ResponseWriter) {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g., Last-Modified might be useful if the
	// response does not have an ETag field).
	h := w.Header()
	delete(h, "Content-Type")
	delete(h, "Content-Length")
	if h.Get("Etag") != "" {
		delete(h, "Last-Modified")
	}
	w.WriteHeader(http.StatusNotModified)
}

var unixEpochTime = time.Unix(0, 0)

// isZeroTime reports whether t is obviously unspecified (either zero or Unix()=0).
func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

type fileServerErrorKind int

const (
	ErrBadRequest fileServerErrorKind = iota
	ErrRedirect
	ErrFsOpen
	ErrFsStat
	ErrDirFound
)

type Err struct {
	httperror.HttpErr
	Kind     fileServerErrorKind
	Pathname string
	Dirstat  os.FileInfo
}

func newErr(statusHint int, kind fileServerErrorKind, message string, pathname string, cause error) *Err {
	var msg *string
	if message != "" {
		msg = &message
	}
	return &Err{
		httperror.HttpErr{
			goerror.CodedErr{
				goerror.GoErr{msg, cause},
				statusHint,
			},
			true,
			nil,
		},
		kind,
		pathname,
		nil,
	}
}

func newErrDirFound(pathname string, dirstat os.FileInfo) *Err {
	e := newErr(http.StatusForbidden, ErrDirFound, "dir found", pathname, nil)
	e.Dirstat = dirstat
	return e
}

func newErrRedirect(r *http.Request, location string, pathname string) *Err {
	path := handlerutils.UnsafeRedirectPath(r, location)
	e := newErr(http.StatusMovedPermanently, ErrRedirect, "", pathname, nil)
	e.Headers().Set("Location", path)
	return e
}
