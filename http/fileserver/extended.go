package fileserver

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/responder"
	"github.com/prasannavl/go-gluons/log"
)

func New(root http.FileSystem) mchain.Handler {
	return &FileServer{root}
}

type FileServer struct {
	root http.FileSystem
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)
		return httperror.NewWithCause(e.StatusHint, "fileserver error", e, false)
	}
	return nil
}

func NewEx(root http.FileSystem, errorTemplatesRoot string) mchain.Handler {
	return &FileServerEx{FileServer{root}, errorTemplatesRoot, ".html"}
}

type FileServerEx struct {
	FileServer
	ErrorTemplatesRoot string
	TemplateSuffix     string
}

func (f *FileServerEx) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)

		log.With("path", r.URL.Path).
			With("status", e.StatusHint).
			With("cause", e.Cause).
			Warnf("fileserver: %v %v", e, e.Data)

		switch e.Kind {
		case ErrFsStat, ErrFsOpen:
			t := errToTemplate(e, f.ErrorTemplatesRoot, f.TemplateSuffix)
			responder.SendWithStatus(w, r,
				httperror.ErrorCode(e.StatusHint), &t)
		case ErrRedirect:
			localRedirect(w, r, e.Data.(string))
		case ErrDirFound:
			return httperror.NewWithCause(e.StatusHint, "fileserver: directory listing disabled", e, false)
			//  Alternatively, handle the dir listing.

			//	f := e.Data.(http.File)
			//	finished := HandleDirPrelude(w, r, f)
			//	if !finished {
			//		// Handle the dir listing here.
			//	}
		}
	}
	return err
}

func statusCodeFromHint(hint int) string {
	return strconv.Itoa(httperror.ErrorCode(hint))
}

func errToTemplate(e *Err, errorTemplatesPath string, templateSuffix string) responder.TemplateFilesContent {
	return responder.TemplateFilesContent{
		Data: e.Error(),
		TemplateFiles: []string{
			filepath.Join(errorTemplatesPath, statusCodeFromHint(e.StatusHint)+templateSuffix)}}
}
