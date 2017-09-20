package log

import "os"

var (
	NopLogger = newNopLogger()
	g         *Logger
)

func init() {
	g = &Logger{
		sink: &StreamSink{
			Stream:    os.Stderr,
			Formatter: DefaultTextFormatterForHuman,
		},
		filter: AllLevelsFilter,
	}
}

func newNopLogger() *Logger {
	l := New(NopSink{})
	SetFlags(l, 0)
	return l
}

func New(sink Sink) *Logger {
	return &Logger{sink, AllLevelsFilter, nil, FlagTime}
}

func SetLogger(l *Logger) {
	if l == nil {
		g = NopLogger
	} else {
		g = l
	}
}

func GetLogger() *Logger {
	return g
}

func GetSink(logger *Logger) Sink {
	return logger.sink
}

func GetFilter(logger *Logger) func(Level) bool {
	return logger.filter
}

func SetFilter(logger *Logger, filter func(Level) bool) {
	if filter == nil {
		filter = AllLevelsFilter
	}
	logger.filter = filter
}

func GetFields(logger *Logger) []Field {
	return logger.fields
}

func GetFlags(logger *Logger) loggerFlags {
	return logger.flags
}

func SetFlags(logger *Logger, flags loggerFlags) {
	logger.flags = flags
}

// Log methods

func Log(lvl Level, message string) {
	logCore(g, lvl, message, nil, skipFramesNum)
}

func Error(message string) {
	logCore(g, ErrorLevel, message, nil, skipFramesNum)
}

func Warn(message string) {
	logCore(g, WarnLevel, message, nil, skipFramesNum)
}

func Info(message string) {
	logCore(g, InfoLevel, message, nil, skipFramesNum)
}

func Debug(message string) {
	logCore(g, DebugLevel, message, nil, skipFramesNum)
}

func Trace(message string) {
	logCore(g, TraceLevel, message, nil, skipFramesNum)
}

// Logv methods

func Logv(lvl Level, args ...interface{}) {
	logCore(g, lvl, "", args, skipFramesNum)
}

func Errorv(args ...interface{}) {
	logCore(g, ErrorLevel, "", args, skipFramesNum)
}

func Warnv(args ...interface{}) {
	logCore(g, WarnLevel, "", args, skipFramesNum)
}

func Infov(args ...interface{}) {
	logCore(g, InfoLevel, "", args, skipFramesNum)
}

func Debugv(args ...interface{}) {
	logCore(g, DebugLevel, "", args, skipFramesNum)
}

func Tracev(args ...interface{}) {
	logCore(g, TraceLevel, "", args, skipFramesNum)
}

// Logf methods

func Logf(lvl Level, format string, args ...interface{}) {
	logCore(g, lvl, format, args, skipFramesNum)
}

func Errorf(format string, args ...interface{}) {
	logCore(g, ErrorLevel, format, args, skipFramesNum)
}

func Warnf(format string, args ...interface{}) {
	logCore(g, WarnLevel, format, args, skipFramesNum)
}

func Infof(format string, args ...interface{}) {
	logCore(g, InfoLevel, format, args, skipFramesNum)
}

func Debugf(format string, args ...interface{}) {
	logCore(g, DebugLevel, format, args, skipFramesNum)
}

func Tracef(format string, args ...interface{}) {
	logCore(g, TraceLevel, format, args, skipFramesNum)
}

func Flush() {
	g.Flush()
}

func With(name string, val interface{}) *Logger {
	return g.With(name, val)
}

func WithFields(fields []Field) *Logger {
	return g.WithFields(fields)
}
