package fileserver

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/prasannavl/goerror/httperror"

	"github.com/prasannavl/go-gluons/httputils/responder"
)

func New(root http.FileSystem, templatesBasePath string) http.Handler {
	return &FileServerEx{root, templatesBasePath, ".html"}
}

type FileServerEx struct {
	root              http.FileSystem
	templatesBasePath string
	templatesSuffix   string
}

func statusCodeFromHint(hint int) string {
	return strconv.Itoa(httperror.ErrorCode(hint))
}

func (f *FileServerEx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ServeRequestPath(w, r, f.root)
	if err != nil {
		e := err.(*Err)
		switch e.Kind {
		case ErrFsStat, ErrFsOpen:
			responder.Send(w, r, &responder.TemplateFilesContent{
				Data: e.Error(),
				TemplateFiles: []string{filepath.Join(f.templatesBasePath,
					statusCodeFromHint(e.StatusHint)+f.templatesSuffix)}})
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
