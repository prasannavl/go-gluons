package middleware

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"

	"github.com/prasannavl/mchain"
)

func PanicRecoveryMiddleware(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) (err error) {
		defer mchain.RecoverIntoError(&err)
		err = next.ServeHTTP(w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func ErrorHandlerMiddleware(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) (err error) {
		err = next.ServeHTTP(w, r)
		if err != nil {
			switch e := err.(type) {
			case httperror.HttpError:
				w.WriteHeader(e.Code())
			case error:
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return err
	}
	return mchain.HandlerFunc(f)
}
