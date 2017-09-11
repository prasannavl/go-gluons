package log

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prasannavl/go-gluons/ansicode"
)

var initTime = time.Now()

type ColorStringer interface {
	ColorString() string
}

func DefaultTextFormatterForHuman(r *Record) string {
	var buf bytes.Buffer
	var timeFormat string
	t := r.Meta.Time
	if t.Sub(initTime).Hours() > 24 {
		timeFormat = "03:04:05 (Jan 02)"
	} else {
		timeFormat = "03:04:05"
	}
	buf.WriteString(t.Format(timeFormat))
	buf.WriteString("  " + PaddedString(GetLogLevelString(r.Meta.Level), 5) + "  ")
	args := r.Args
	if r.Format == "" {
		fmt.Fprint(&buf, args...)
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, r.Format, args...)
	} else {
		buf.WriteString(r.Format)
	}
	for _, x := range r.Fields() {
		fmt.Fprintf(&buf, " %s=%v ", x.Name, x.Value)
	}
	buf.WriteString("\r\n")
	return buf.String()
}

func DefaultColorTextFormatterForHuman(r *Record) string {
	var buf bytes.Buffer
	t := r.Meta.Time
	var timeFormat string
	if t.Sub(initTime).Hours() > 24 {
		timeFormat = "03:04:05 (Jan 02)"
	} else {
		timeFormat = "03:04:05"
	}
	buf.WriteString(ansicode.BlackBright + t.Format(timeFormat) + ansicode.Reset)
	buf.WriteString("  " + logLevelColoredString(r.Meta.Level) + "  ")
	args := r.Args
	for i, a := range args {
		if colorable, ok := a.(ColorStringer); ok {
			args[i] = colorable.ColorString()
		}
	}
	if r.Format == "" {
		fmt.Fprint(&buf, args...)
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, r.Format, args...)
	} else {
		buf.WriteString(r.Format)
	}
	for _, x := range r.Fields() {
		value := x.Value
		if colorable, ok := value.(ColorStringer); ok {
			value = colorable.ColorString()
		}
		fmt.Fprintf(&buf, " %s=%v ", HashColoredText(x.Name), value)
	}
	buf.WriteString("\r\n")
	return buf.String()
}

var colorMap = []string{
	ansicode.BlackBright,
	ansicode.Cyan,
	ansicode.Green,
	ansicode.Magenta,
}

var colorMapLen = len(colorMap)

func HashColoredText(name string) string {
	const maxIterations = 10
	l := len(name)
	if l > 10 {
		l = 10
	}
	for i, x := range name {
		if i > maxIterations {
			break
		}
		l += int(x)
	}
	index := l % colorMapLen
	if index < 0 {
		index = 0
	}
	return colorMap[index] + name + ansicode.Reset
}

func DefaultTextFormatter(r *Record) string {
	var buf bytes.Buffer
	buf.WriteString(r.Meta.Time.Format(time.RFC3339))
	buf.WriteString("," + GetLogLevelString(r.Meta.Level) + ",")
	args := r.Args
	if r.Format == "" {
		buf.WriteString(strconv.Quote(fmt.Sprint(args...)))
	} else if len(args) > 0 {
		buf.WriteString(strconv.Quote(fmt.Sprintf(r.Format, args...)))
	} else {
		buf.WriteString(strconv.Quote(r.Format))
	}
	ctx := r.Fields()
	for _, x := range ctx {
		buf.WriteString("," + strconv.Quote(x.Name) + "=" + strconv.Quote(fmt.Sprint(x.Value)))
	}
	buf.WriteString("\r\n")
	return buf.String()
}

func PaddedString(s string, width int) string {
	diff := width - len(s)
	if diff == 1 {
		return s + " "
	} else if diff > 1 {
		return s + strings.Repeat(" ", diff)
	}
	return s
}

func GetLogLevelString(lvl Level) string {
	switch lvl {
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	}
	return "msg"
}

func logLevelColoredString(lvl Level) string {
	return GetLogLevelColoredMsg(lvl, PaddedString(GetLogLevelString(lvl), 5))
}

func GetLogLevelColoredMsg(lvl Level, msg string) string {
	switch lvl {
	case InfoLevel:
		return ansicode.Blue + msg + ansicode.Reset
	case WarnLevel:
		return ansicode.YellowBright + msg + ansicode.Reset
	case ErrorLevel:
		return ansicode.RedBright + msg + ansicode.Reset
	case DebugLevel:
		return ansicode.White + msg + ansicode.Reset
	case TraceLevel:
		return ansicode.BlackBright + msg + ansicode.Reset
	}
	return ansicode.BlackBright + msg + ansicode.Reset
}
