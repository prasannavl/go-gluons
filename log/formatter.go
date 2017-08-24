package log

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
)

func DefaultTextFormatterForHuman(r *Record) string {
	var buf bytes.Buffer
	buf.WriteString(r.Meta.Time.Format("03:04:05 (Jan 02)"))
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
	buf.WriteString(color.HiBlackString(r.Meta.Time.Format("03:04:05 (Jan 02)")))
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
	return ""
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
	return " MSG  "
}

func GetLogLevelColoredString(lvl Level) string {
	return GetLogLevelColoredMsg(lvl, GetLogLevelString(lvl))
}

func GetLogLevelColoredMsg(lvl Level, msg string) string {
	switch lvl {
	case InfoLevel:
		return color.HiBlueString(msg)
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
