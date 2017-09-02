package reqcontext

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	mx "github.com/go-chi/chi/middleware"
	"github.com/prasannavl/go-gluons/ansicode"
	"github.com/prasannavl/go-gluons/httputils/app/responder"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/goerror/errutils"
)

func CreateLogHandler(requestLogLevel log.Level) middleware {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			ww := w.(mx.WrapResponseWriter)
			defer func() {
				reqctx := FromRequest(r)
				logger := &reqctx.Logger
				if err := recover(); err != nil {
					responder.SendErrorText(w, err)
					stack := debug.Stack()
					LogErrorStack(logger, err, stack)
				}
				sizeRaw := ww.BytesWritten()
				if sizeRaw > 0 {
					logger.Logf(
						requestLogLevel,
						"%s %v %v %s %s",
						r.Method,
						ColoredHttpStatus(ww.Status()),
						ColoredDuration(time.Since(startTime)),
						sizeString(sizeRaw),
						r.URL.String())
				} else {
					logger.Logf(
						requestLogLevel,
						"%s %v %v %s",
						r.Method,
						ColoredHttpStatus(ww.Status()),
						ColoredDuration(time.Since(startTime)),
						r.URL.String())
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
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
		logger.Errorf("[[stack]]\r\n%s[[stack:end]]\r\n", stack)
	}
}

const (
	sizeKB = 1 * 1000
	sizeMB = sizeKB * 1000
	sizeGB = sizeMB * 1000
)

func sizeString(size int) string {
	s := strconv.Itoa(size)
	if size < sizeKB {
		return s + "B"
	}
	if size < sizeMB {
		return s + "KB"
	}
	if size < sizeGB {
		return s + "MB"
	}
	return s + "GB"
}

type ColoredHttpStatus int

func (self ColoredHttpStatus) ColorString() string {
	return coloredStatusIntString(int(self))
}

var statusColorMap = []string{
	ansicode.Cyan,
	ansicode.Green,
	ansicode.Magenta,
	ansicode.Yellow,
	ansicode.RedBright,
}

var statusColorMapLen = len(statusColorMap)

func coloredStatusIntString(code int) string {
	index := code
	if index < 100 {
		return statusColorMap[0]
	}
	index = index / 100
	index = (index - 1) % statusColorMapLen
	return statusColorMap[index] + strconv.Itoa(code) + ansicode.Reset
}

type ColoredDuration time.Duration

func (self ColoredDuration) ColorString() string {
	return coloredTimeString(time.Duration(self))
}

var timeColorMap = []string{
	ansicode.BlackBright,
	ansicode.Yellow,
	ansicode.RedBright,
	ansicode.Red,
	ansicode.RedBrightBg + ansicode.White,
	ansicode.RedBg + ansicode.White + ansicode.Bold,
}

var timeColorMapLen = len(statusColorMap)

func coloredTimeString(duration time.Duration) string {
	index := 0
	t := duration.Nanoseconds() / 1000
	millis := t / 1000
	if millis > 80 {
		index++
		if millis > 200 {
			index++
			if millis > 500 {
				index++
				if millis > 1000 {
					index++
					if millis > 2000 {
						index++
					}
				}
			}
		}
	}
	return timeColorMap[index] + fmt.Sprintf("%v", duration) + ansicode.Reset
}
