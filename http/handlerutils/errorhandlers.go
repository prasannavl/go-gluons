package handlerutils

import (
	"net/http"

	"github.com/prasannavl/go-errors/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/reqcontext"
	"github.com/prasannavl/go-gluons/http/writer"
)

func StatusErrorHandler(status int, logErrors bool) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		if logErrors {
			logger := reqcontext.GetRequestLogger(r)
			logger.Errorf("error-handler: %v", err)
		}
		if ww, ok := w.(writer.ResponseWriter); ok {
			if !ww.IsStatusWritten() {
				w.WriteHeader(status)
			}
			return
		}
		w.WriteHeader(status)
	}
}

func HttpErrorHandler(fallbackStatus int, logErrors bool) mchain.ErrorHandler {
	statusErrorHandler := StatusErrorHandler(fallbackStatus, logErrors)
	return func(err error, w http.ResponseWriter, r *http.Request) {
		switch e := err.(type) {
		case httperror.HttpError:
			if httperror.IsServerErrorCode(e.Code()) && logErrors {
				logger := reqcontext.GetRequestLogger(r)
				logger.Errorf("error-handler: %v", e)
			}
			if ww, ok := w.(writer.ResponseWriter); ok {
				if !ww.IsStatusWritten() {
					w.WriteHeader(e.Code())
					e.Headers().Write(w)
				}
				return
			}
			w.WriteHeader(e.Code())
			e.Headers().Write(w)
		case error:
			statusErrorHandler(err, w, r)
		}
	}
}
