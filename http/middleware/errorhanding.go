package middleware

import (
	"net/http"

	"github.com/prasannavl/go-gluons/http/handlerutils"

	"github.com/prasannavl/go-gluons/http/writer"
	"github.com/prasannavl/mchain"
)

func ErrorHandlerMiddleware(next mchain.Handler) mchain.Handler {
	handler := handlerutils.HttpErrorHandler(http.StatusInternalServerError, false)
	f := func(w http.ResponseWriter, r *http.Request) (err error) {
		err = next.ServeHTTP(w, r)
		ww := w.(writer.ResponseWriter)
		if ww.IsHijacked() {
			return err
		}
		if err != nil {
			handler(err, w, r)
		}
		return err
	}
	return mchain.HandlerFunc(f)
}

func PanicRecoveryMiddleware(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) (err error) {
		defer mchain.RecoverIntoError(&err)
		err = next.ServeHTTP(w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}
