package middleware

import (
	"log"
	"net/http"
	"pvl/api-core/reqcontext"
	"time"

	"github.com/prasannavl/mchain"
)

func RequestContextInitHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		return next.ServeHTTP(w, reqcontext.WithRequestContext(r, &reqcontext.RequestContext{}))
	}
	return mchain.HandlerFunc(f)
}

func RequestLogHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		err := next.ServeHTTP(w, r)
		c := reqcontext.FromRequest(r)
		log.Printf("[%[1]s] %[3]v %[4]s |%[2]s|", r.Method, c.RequestID, c.EndTime.Sub(c.StartTime), r.URL.String())
		return err
	}
	return mchain.HandlerFunc(f)
}

func RequestDurationHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		c := reqcontext.FromRequest(r)
		c.StartTime = time.Now()
		err := next.ServeHTTP(w, r)
		c.EndTime = time.Now()
		return err
	}
	return mchain.HandlerFunc(f)
}
