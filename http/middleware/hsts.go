package middleware

import (
	"net/http"
	"strconv"

	"github.com/prasannavl/mchain"
)

//  Ref: (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security)
//
// 	Syntax:
//		Strict-Transport-Security: max-age=<expire-time>
//		Strict-Transport-Security: max-age=<expire-time>; includeSubDomains
//		Strict-Transport-Security: max-age=<expire-time>; preload
//

func StrictTransportSecMiddleware(maxAgeSecs int, includeSubdomains bool, preload bool) mchain.Middleware {
	return func(next mchain.Handler) mchain.Handler {

		val := strconv.Itoa(maxAgeSecs)
		if includeSubdomains {
			val += "; includeSubDomains"
		}
		if preload {
			val += "; preload"
		}

		f := func(w http.ResponseWriter, r *http.Request) (err error) {
			w.Header().Set("Strict-Transport-Security", val)
			err = next.ServeHTTP(w, r)
			return err
		}
		return mchain.HandlerFunc(f)
	}
}
