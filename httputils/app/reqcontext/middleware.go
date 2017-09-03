package reqcontext

import (
	"net/http"

	chimx "github.com/go-chi/chi/middleware"

	"runtime/debug"

	"github.com/prasannavl/go-gluons/log"
)

func CreateInitMiddleware(l *log.Logger) middleware {
	m := func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			ww := chimx.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, WithRequestContext(r, &RequestContext{
				Logger: *l,
			}))
		}
		return http.HandlerFunc(f)
	}
	return m
}

func CreateRecoveryMiddleware(
	errorHandler func(err interface{}, r *http.Request),
	responder func(w http.ResponseWriter, r *http.Request)) middleware {
	m := func(next http.Handler) http.Handler {
		if errorHandler == nil {
			errorHandler = DefaultRecoveryErrorHandler
		}
		if responder == nil {
			responder = DefaultRecoveryResponder
		}
		f := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					errorHandler(err, r)
					responder(w, r)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
	return m
}

func DefaultRecoveryErrorHandler(err interface{}, r *http.Request) {
	ctx := FromRequest(r)
	ctx.Recovery.Error = err
	ctx.Recovery.Stack = debug.Stack()
}

func DefaultRecoveryResponder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}
