package reqcontext

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	mx "github.com/go-chi/chi/middleware"
	"github.com/prasannavl/go-gluons/ansicode"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/goerror/errutils"
)

func CreateLogMiddleware(requestLogLevel log.Level) middleware {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			ww := w.(mx.WrapResponseWriter)
			startTime := time.Now()

			next.ServeHTTP(w, r)

			ctx := FromRequest(r)
			logger := &ctx.Logger
			if err := ctx.Recovery.Error; err != nil {
				LogErrorStack(logger, err, ctx.Recovery.Stack)
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
	ansicode.YellowBright,
	ansicode.RedBright,
}

var sizeColorsLen = len(statusColors)

const (
	sizeKB = 1 * 1000
)

var sizeColorLimits = [...]int{80 * sizeKB, 400 * sizeKB}

func (c ColoredTransferSize) ColorString() string {
	str := sizeString(int(c))
	index := 0
	current := int(c)
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

var sizeSuffixes = [...]string{"b", "kb", "mb", "gb"}

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

func (c ColoredHttpStatus) ColorString() string {
	return coloredStatusIntString(int(c))
}

var statusColors = [...]string{
	ansicode.CyanBright,
	ansicode.Green,
	ansicode.MagentaBright,
	ansicode.YellowBright,
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

func (c ColoredDuration) ColorString() string {
	return coloredTimeString(time.Duration(c))
}

var timeColors = [...]string{
	ansicode.BlackBright,
	ansicode.YellowBright,
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
