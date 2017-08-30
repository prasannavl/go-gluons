package middleware

import (
	"net/http"

	"github.com/prasannavl/mchain"
)

func RecoverPanicHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) (err error) {
		defer mchain.RecoverIntoError(&err)
		err = next.ServeHTTP(w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}
