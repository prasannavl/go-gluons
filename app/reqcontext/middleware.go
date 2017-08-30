package reqcontext

import (
	"net/http"
	"time"

	"github.com/prasannavl/goerror/httperror"

	"github.com/prasannavl/go-grab/log"
	"github.com/prasannavl/go-starter-api/app/responder"
	"github.com/prasannavl/goerror/errutils"

	"github.com/prasannavl/mchain"
)

func CreateInitHandler(l *log.Logger) mchain.Middleware {
	m := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			return next.ServeHTTP(w, WithRequestContext(r, &RequestContext{
				Logger: *l,
			}))
		}
		return mchain.HandlerFunc(f)
	}
	return m
}

func ErrorHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		err := next.ServeHTTP(w, r)
		if err != nil {
			c := FromRequest(r)
			iter := errutils.MakeIteratorLimited(err, 10)
			var httpErr httperror.Error
			for {
				e := iter.Next()
				if e == nil {
					break
				}
				if herr, ok := e.(httperror.Error); ok {
					httpErr = herr
				}
				c.Logger.Errorf("cause: %s => %#v ", e.Error(), e)
			}
			if httpErr != nil {
				responder.SendHttpError(httpErr, w, r)
			} else {
				responder.SendStatus(http.StatusInternalServerError, w)
			}
		}
		return nil
	}
	return mchain.HandlerFunc(f)
}

func LogHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		c := FromRequest(r)
		err := next.ServeHTTP(w, r)
		if err != nil {
			c.Logger.With("err", err.Error()).
				Errorf("%s %v %s", r.Method, c.EndTime.Sub(c.StartTime), r.URL.String())
		} else {
			c.Logger.Tracef("%s %v %s", r.Method, c.EndTime.Sub(c.StartTime), r.URL.String())
		}
		return err
	}
	return mchain.HandlerFunc(f)
}

func DurationHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		c := FromRequest(r)
		c.StartTime = time.Now()
		err := next.ServeHTTP(w, r)
		c.EndTime = time.Now()
		return err
	}
	return mchain.HandlerFunc(f)
}
