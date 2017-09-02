package reqcontext

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/prasannavl/go-gluons/httputils/app/responder"

	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/goerror/errutils"
)

func LogHandler(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			reqctx := FromRequest(r)
			logger := &reqctx.Logger
			if err := recover(); err != nil {
				responder.SendErrorText(w, err)
				stack := debug.Stack()
				LogErrorStack(logger, err, stack)
			}
			logger.Tracef("%s %v %s", r.Method, time.Since(startTime), r.URL.String())
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func LogError(logger *log.Logger, e interface{}) {
	if err, ok := e.(error); ok {
		iter := errutils.MakeIteratorLimited(err, 10)
		for {
			e := iter.Next()
			if e == nil {
				break
			}
			logger.Errorf("cause: %s => %#v ", e.Error(), e)
		}
	} else {
		logger.Errorf("%#v", e)
	}
}

func LogErrorStack(logger *log.Logger, err interface{}, stack []byte) {
	LogError(logger, err)
	if len(stack) > 0 {
		logger.Errorf("stack:\r\n%s\r\n", stack)
	}
}
