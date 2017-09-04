package fileserver

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func FileServer(root http.FileSystem) http.Handler {
	return &fileHandler{root}
}

type fileHandler struct {
	root http.FileSystem
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)
		switch e.Kind {
		case ErrBadRequest, ErrFsStat, ErrFsOpen:
			http.Error(w, e.Error(), e.StatusHint)
		case ErrRedirect:
			LocalRedirect(w, r, e.Data.(string))
		case ErrDirFound:
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			//  Alternatively, handle the dir listing.

			//	f := e.Data.(http.File)
			//	finished := HandleDirPrelude(w, r, f)
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
	if containsDotDot(r.URL.Path) {
		// Too many programs use r.URL.Path to construct the argument to
		// serveFile. Reject the request under the assumption that happened
		// here and ".." may not be wanted.
		// Note that name might not contain "..", for example if code (still
		// incorrectly) used filepath.Join(myDir, r.URL.Path).
		return &Err{
			Kind:       ErrBadRequest,
			StatusHint: http.StatusBadRequest,
			Data:       "invalid URL path",
		}
	}
	dir, file := filepath.Split(name)
	return serveFile(w, r, http.Dir(dir), file, false)
}

func containsDotDot(v string) bool {
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

	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, indexPage) {
		return &Err{
			Kind: ErrRedirect, Data: "./",
			StatusHint: http.StatusMovedPermanently,
		}
	}

	f, err := fs.Open(name)
	if err != nil {
		return &Err{
			Kind:       ErrFsOpen,
			StatusHint: http.StatusNotFound,
			Cause:      err,
		}
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return &Err{
			Kind:       ErrFsStat,
			StatusHint: http.StatusInternalServerError,
			Cause:      err,
		}
	}

	if redirect {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path
		if d.IsDir() {
			if url[len(url)-1] != '/' {
				return &Err{
					Kind: ErrRedirect, Data: path.Base(url) + "/",
					StatusHint: http.StatusMovedPermanently,
				}
			}
		} else {
			if url[len(url)-1] == '/' {
				return &Err{
					Kind: ErrRedirect, Data: "../" + path.Base(url),
					StatusHint: http.StatusMovedPermanently,
				}
			}
		}
	}

	// redirect if the directory name doesn't end in a slash
	if d.IsDir() {
		url := r.URL.Path
		if url[len(url)-1] != '/' {
			return &Err{
				Kind: ErrRedirect, Data: path.Base(url) + "/",
				StatusHint: http.StatusMovedPermanently,
			}
		}
	}

	// use contents of index.html for directory, if present
	if d.IsDir() {
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := fs.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				name = index
				d = dd
				f = ff
			}
		}
	}

	// Still a directory? (we didn't find an index.html file)
	if d.IsDir() {
		return &Err{
			Kind:       ErrDirFound,
			Data:       f,
			StatusHint: 0,
		}
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
	return nil
}

// LocalRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func LocalRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

func HandleDirPrelude(w http.ResponseWriter, r *http.Request, dir http.File) (finished bool) {
	// Note: the error is ignored for now, since this already would have been executed
	// in before. If there's an issue here, something has gone seriously wrong,
	// in which case, we panic anyway.
	d, _ := dir.Stat()
	if checkIfModifiedSince(r, d.ModTime()) == condFalse {
		writeNotModified(w)
		return true
	}
	w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
	return false
}

// condResult is the result of an HTTP request precondition check.
// See https://tools.ietf.org/html/rfc7232 section 3.
type condResult int

const (
	condNone condResult = iota
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
	Kind       fileServerErrorKind
	StatusHint int
	Data       interface{}
	Cause      error
}

func (f *Err) Error() string {
	if f.StatusHint == 0 {
		return "application action required"
	}
	return http.StatusText(f.StatusHint)
}
