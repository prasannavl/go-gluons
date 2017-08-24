package log

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
)

var initTime = time.Now()

func DefaultTextFormatterForHuman(r *Record) string {
	var buf bytes.Buffer
	var timeFormat string
	t := r.Meta.Time
	if t.Sub(initTime).Hours() > 24 {
		timeFormat = "03:04:05 (Jan 02)"
	} else {
		timeFormat = "03:04:05"
	}
	buf.WriteString(color.HiBlackString(t.Format(timeFormat)))
	buf.WriteString("  " + GetLogLevelString(r.Meta.Level) + "  ")
	args := r.Args
	if r.Format == "" {
		fmt.Fprint(&buf, args...)
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, r.Format, args...)
	} else {
		buf.WriteString(r.Format)
	}
	for _, x := range r.Context() {
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
	buf.WriteString(color.HiBlackString(t.Format(timeFormat)))
	buf.WriteString("  " + GetLogLevelColoredString(r.Meta.Level) + "  ")
	args := r.Args
	if r.Format == "" {
		fmt.Fprint(&buf, args...)
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, r.Format, args...)
	} else {
		buf.WriteString(r.Format)
	}
	for _, x := range r.Context() {
		fmt.Fprintf(&buf, " %s=%v ", x.Name, x.Value)
	}
	buf.WriteString("\r\n")
	return buf.String()
}

func DefaultTextFormatter(r *Record) string {
	var buf bytes.Buffer
	buf.WriteString(r.Meta.Time.Format(time.RFC3339))
	buf.WriteString("," + strconv.Itoa(int(r.Meta.Level)) + ",")
	args := r.Args
	if r.Format == "" {
		fmt.Fprintf(&buf, "%q", fmt.Sprint(args...))
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, "%q", fmt.Sprintf(r.Format, args...))
	} else {
		fmt.Fprintf(&buf, "%q", r.Format)
	}
	ctx := r.Context()
	for _, x := range ctx {
		fmt.Fprintf(&buf, ",%q=%q", x.Name, x.Value)
	}
	buf.WriteString("\r\n")
	return buf.String()
}

func GetLogLevelString(lvl Level) string {
	switch lvl {
	case InfoLevel:
		return "info "
	case WarnLevel:
		return "warn "
	case ErrorLevel:
		return "error"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	}
	return " msg  "
}

func GetLogLevelColoredString(lvl Level) string {
	return GetLogLevelColoredMsg(lvl, GetLogLevelString(lvl))
}

func GetLogLevelColoredMsg(lvl Level, msg string) string {
	switch lvl {
	case InfoLevel:
		return color.BlueString(msg)
	case WarnLevel:
		return color.HiYellowString(msg)
	case ErrorLevel:
		return color.HiRedString(msg)
	case DebugLevel:
		return color.WhiteString(msg)
	case TraceLevel:
		return color.HiBlackString(msg)
	}
	return color.HiBlackString(msg)
}
