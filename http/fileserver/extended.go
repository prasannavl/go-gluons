package fileserver

import (
	"net/http"
	"strconv"

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

func NewEx(root http.FileSystem, notFoundHandler mchain.HandlerFunc) mchain.Handler {
	return &FileServerEx{FileServer{root}, true, notFoundHandler}
}

type FileServerEx struct {
	FileServer
	RedirectsEnabled bool
	NotFoundHandler  mchain.HandlerFunc
}

func (f *FileServerEx) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)

		log.With("path", r.URL.Path).
			With("status", e.Code()).
			With("cause", e.Cause()).
			Warnf("fileserver: %v %v", e, e.Data)

		switch e.Kind {
		case ErrRedirect:
			if f.RedirectsEnabled {
				localRedirect(w, r, e.Data.(string))
				return nil
			}
		case ErrFsOpen:
			if f.NotFoundHandler != nil {
				return f.NotFoundHandler(w, r)
			}
		}
	}
	return err
}

func statusCodeFromHint(hint int) string {
	return strconv.Itoa(httperror.ErrorCode(hint))
}
