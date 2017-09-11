package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prasannavl/go-gluons/ansicode"
	"github.com/prasannavl/go-gluons/http/writer"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/goerror/errutils"
	"github.com/prasannavl/mchain"
)

func CreateLogMiddleware(requestLogLevel log.Level) mchain.Middleware {
	return func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			ww := w.(writer.ResponseWriter)
			startTime := time.Now()
			err := next.ServeHTTP(w, r)
			ctx := FromRequest(r)
			if err != nil {
				LogError(&ctx.Logger, err)
				LogErrorStack(&ctx.Logger, ctx.ErrorStacks...)
			}
			logger := &ctx.Logger
			sizeRaw := ww.BytesWritten()
			if sizeRaw > 0 {
				logger.Logf(
					requestLogLevel,
					"%s %v %v %v %s %s",
					r.Method,
					ColoredHttpStatus(ww.Status()),
					ColoredDuration(time.Since(startTime)),
					ColoredTransferSize(sizeRaw),
					r.RequestURI,
					r.RemoteAddr,
				)
			} else {
				logger.Logf(
					requestLogLevel,
					"%s %v %v %s %s",
					r.Method,
					ColoredHttpStatus(ww.Status()),
					ColoredDuration(time.Since(startTime)),
					r.RequestURI,
					r.RemoteAddr)
			}
			return nil
		}
		return mchain.HandlerFunc(f)
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
			if errutils.HasMessage(e) {
				logger.Errorf("cause: %s => %#v ", e.Error(), e)
			}
		}
	} else {
		logger.Errorf("%#v", e)
	}
}

func LogErrorStack(logger *log.Logger, stacks ...[]byte) {
	var i = 0
	for _, stack := range stacks {
		if len(stack) > 0 {
			logger.Errorf("[[stack-%d]]\r\n%s[[stack:end]]\r\n", i, stack)
			i++
		}
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
