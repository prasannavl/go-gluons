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
						ColoredTransferSize(sizeRaw),
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

type ColoredTransferSize int

var sizeColors = [...]string{
	ansicode.BlackBright,
	ansicode.Yellow,
	ansicode.Red,
}

var sizeColorsLen = len(statusColors)

const (
	sizeKB = 1 * 1000
)

var sizeColorLimits = [...]int{80 * sizeKB, 400 * sizeKB}

func (self ColoredTransferSize) ColorString() string {
	str := sizeString(int(self))
	index := 0
	current := int(self)
	for _, limit := range sizeColorLimits {
		if current > limit {
			index++
		} else {
			break
		}
	}
	if index >= sizeColorsLen {
		index = sizeColorsLen - 1
	}
	return sizeColors[index] + str + ansicode.Reset
}

var sizeSuffixes = [...]string{"B", "KB", "MB", "GB"}

func sizeString(size int) string {
	index := 0
	i := 1000
	for ; ; i *= 1000 {
		if size > i {
			index++
		} else {
			break
		}
	}
	maxLen := len(sizeSuffixes)
	if index >= maxLen {
		size = maxLen - 1
	}
	n := float64(size) / float64(i/1000)
	s := strconv.FormatFloat(n, 'f', -1, 64)
	return s + sizeSuffixes[index]
}

type ColoredHttpStatus int

func (self ColoredHttpStatus) ColorString() string {
	return coloredStatusIntString(int(self))
}

var statusColors = [...]string{
	ansicode.Cyan,
	ansicode.Green,
	ansicode.Magenta,
	ansicode.Yellow,
	ansicode.RedBright,
}

var statusColorsLen = len(statusColors)

func coloredStatusIntString(code int) string {
	index := code
	if index < 100 {
		return statusColors[0]
	}
	index = index / 100
	index = (index - 1) % statusColorsLen
	return statusColors[index] + strconv.Itoa(code) + ansicode.Reset
}

type ColoredDuration time.Duration

func (self ColoredDuration) ColorString() string {
	return coloredTimeString(time.Duration(self))
}

var timeColors = [...]string{
	ansicode.BlackBright,
	ansicode.Yellow,
	ansicode.RedBright,
	ansicode.Red,
	ansicode.RedBrightBg + ansicode.White,
	ansicode.RedBg + ansicode.White + ansicode.Bold,
}

var timeColorsLen = len(statusColors)
var timeColorLimits = [...]int64{80, 200, 500, 1000, 2000}

func coloredTimeString(duration time.Duration) string {
	index := 0
	t := duration.Nanoseconds() / 1000
	millis := t / 1000
	for _, limit := range timeColorLimits {
		if millis > limit {
			index++
		} else {
			break
		}
	}
	if index >= timeColorsLen {
		index = timeColorsLen - 1
	}
	return timeColors[index] + fmt.Sprintf("%v", duration) + ansicode.Reset
}
