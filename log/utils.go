package log

import "strings"

func DisabledFilter(lvl Level) bool {
	return false
}

func AllLevelsFilter(lvl Level) bool {
	return true
}

func ErrorLevelFilter(lvl Level) bool {
	return lvl <= ErrorLevel
}

func WarnLevelFilter(lvl Level) bool {
	return lvl <= WarnLevel
}

func InfoLevelFilter(lvl Level) bool {
	return lvl <= InfoLevel
}

func DebugLevelFilter(lvl Level) bool {
	return lvl <= DebugLevel
}

func TraceLevelFilter(lvl Level) bool {
	return AllLevelsFilter(lvl)
}

func LogLevelFromString(level string) Level {
	switch level {
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "trace", "all":
		return TraceLevel
	case "off":
		return DisabledLevel
	default:
		return Level(^uint(0))
	}
}

func LogLevelString(lvl Level) string {
	switch lvl {
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	case DisabledLevel:
		return "off"
	}
	return "msg"
}

func LogFilterForLevel(lvl Level) func(Level) bool {
	switch lvl {
	case ErrorLevel:
		return ErrorLevelFilter
	case WarnLevel:
		return WarnLevelFilter
	case InfoLevel:
		return InfoLevelFilter
	case DebugLevel:
		return DebugLevelFilter
	case TraceLevel:
		return TraceLevelFilter
	case DisabledLevel:
		return DisabledFilter
	default:
		return InfoLevelFilter
	}
}

func IsValidLevel(lvl Level) bool {
	switch lvl {
	case DisabledLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel, TraceLevel:
		return true
	default:
		return false
	}
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
