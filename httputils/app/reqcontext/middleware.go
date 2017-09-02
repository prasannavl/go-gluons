package reqcontext

import (
	"net/http"

	chimx "github.com/go-chi/chi/middleware"

	"github.com/prasannavl/go-gluons/log"
)

func CreateInitHandler(l *log.Logger) middleware {
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
