package reqcontext

import (
	"net/http"

	"github.com/prasannavl/go-gluons/http/writer"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/mchain"
)

func CreateInitMiddleware(l *log.Logger) mchain.Middleware {
	m := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			ww := writer.NewResponseWriter(w, r.ProtoMajor)
			return next.ServeHTTP(ww, WithRequestContext(r, &RequestContext{
				Logger: *l,
			}))
		}
		return mchain.HandlerFunc(f)
	}
	return m
}
