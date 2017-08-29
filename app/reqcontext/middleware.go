package reqcontext

import (
	"net/http"
	"time"

	"github.com/prasannavl/go-grab/log"
	"github.com/prasannavl/mchain"
)

func RequestContextInitHandler(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
	return next.ServeHTTP(w, WithRequestContext(r, &RequestContext{}))
}

func CreateRequestLogHandler(logger *log.Logger) mchain.SimpleMiddleware {
	f := func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
		c := FromRequest(r)
		c.Logger = *logger
		err := next.ServeHTTP(w, r)
		c.Logger.Tracef("%s %v %s", r.Method, c.EndTime.Sub(c.StartTime), r.URL.String())
		return err
	}
	return mchain.SimpleMiddleware(f)
}

func RequestDurationHandler(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
	c := FromRequest(r)
	c.StartTime = time.Now()
	err := next.ServeHTTP(w, r)
	c.EndTime = time.Now()
	return err
}
