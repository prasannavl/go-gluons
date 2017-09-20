package handlerutils

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/go-gluons/http/writer"
	"github.com/prasannavl/go-gluons/log"
)

func StatusErrorHandler(status int, logErrors bool) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		if logErrors {
			log.Errorf("error-handler: %v", err)
		}
		ww := w.(writer.ResponseWriter)
		if !ww.IsStatusWritten() {
			w.WriteHeader(status)
		}
	}
}

func HttpErrorHandler(fallbackStatus int, logErrors bool) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		ww := w.(writer.ResponseWriter)
		var logger *log.Logger
		ctx := middleware.FromRequest(r)
		if ctx != nil {
			logger = &ctx.Logger
		}
		if logger == nil {
			logger = log.GetLogger()
		}
		switch e := err.(type) {
		case httperror.HttpError:
			if httperror.IsServerErrorCode(e.Code()) && logErrors {
				log.Errorf("error-handler: %v", e)
			}
			if !ww.IsStatusWritten() {
				w.WriteHeader(e.Code())
				e.Headers().Write(w)
			}
		default:
			if logErrors {
				log.Errorf("error-handler: %v", err)
			}
			if !ww.IsStatusWritten() {
				ww.WriteHeader(fallbackStatus)
			}
		}
	}
}