package fileserver

import (
	"net/http"
	"strconv"

	"github.com/prasannavl/goerror/errutils"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/log"
)

func New(root http.FileSystem) mchain.Handler {
	return &FileServer{root}
}

type FileServer struct {
	root http.FileSystem
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return ServeRequestPath(w, r, f.root)
}

func NewEx(root http.FileSystem, notFoundHandler mchain.Handler) mchain.Handler {
	return &FileServerEx{FileServer{root}, true, notFoundHandler}
}

type FileServerEx struct {
	FileServer
	RedirectsEnabled bool
	NotFoundHandler  mchain.Handler
}

func (f *FileServerEx) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)

		l := log.With("path", r.URL.Path).
			With("status", e.Code())
		if !errutils.HasMessage(e) {
			l = l.With("cause", e.Cause())
		}
		switch e.Kind {
		case ErrRedirect:
			if f.RedirectsEnabled {
				e.Headers().Write(w)
				w.WriteHeader(http.StatusMovedPermanently)
				return nil
			}
		case ErrFsOpen:
			l.Warnf("fileserver: kind: %d, %v", e.Kind, e)
			if f.NotFoundHandler != nil {
				return f.NotFoundHandler.ServeHTTP(w, r)
			}
		default:
			l.Warnf("fileserver: kind: %d, %v", e.Kind, e)
		}
	}
	return err
}

func statusCodeFromHint(hint int) string {
	return strconv.Itoa(httperror.ErrorCode(hint))
}
