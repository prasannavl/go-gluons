package log

import (
	"runtime"
	"time"
)

type Level uint

// Power of 2
const (
	DisabledLevel Level = 0
	ErrorLevel          = 1 << (iota - 1)
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type loggerFlags int

const (
	_           loggerFlags = 0
	FlagTime                = 1 << (iota - 1)
	FlagSrcHint
)

type Logger struct {
	sink   Sink
	filter func(Level) bool
	fields []Field
	flags  loggerFlags
}

type Record struct {
	Meta   Metadata
	Format string
	Args   []interface{}
}

type Field struct {
	Name  string
	Value interface{}
}

type Metadata struct {
	Logger *Logger
	Level  Level
	Time   time.Time
	File   string
	Line   int
}

const skipFramesNum = 3

func logCore(l *Logger, lvl Level, format string, args []interface{}, skipStackFramesNum int) {
	if l.IsEnabled(lvl) {
		r := Record{
			Meta:   newMetadata(l, lvl, skipStackFramesNum),
			Format: format,
			Args:   args,
		}
		l.sink.Log(&r)
	}
}

func newMetadata(l *Logger, lvl Level, skip int) Metadata {
	m := Metadata{Logger: l, Level: lvl}
	f := l.flags
	if f&FlagTime == FlagTime {
		m.Time = time.Now()
	}
	if f&FlagSrcHint == FlagSrcHint {
		if _, file, line, ok := runtime.Caller(skip); ok {
			m.File = file
			m.Line = line
		}
	}
	return m
}

// Log methods

func (l *Logger) Log(lvl Level, message string) {
	logCore(l, lvl, message, nil, skipFramesNum)
}

func (l *Logger) Error(message string) {
	logCore(l, ErrorLevel, message, nil, skipFramesNum)
}

func (l *Logger) Warn(message string) {
	logCore(l, WarnLevel, message, nil, skipFramesNum)
}

func (l *Logger) Info(message string) {
	logCore(l, InfoLevel, message, nil, skipFramesNum)
}

func (l *Logger) Debug(message string) {
	logCore(l, DebugLevel, message, nil, skipFramesNum)
}

func (l *Logger) Trace(message string) {
	logCore(l, TraceLevel, message, nil, skipFramesNum)
}

// Logv methods

func (l *Logger) Logv(lvl Level, args ...interface{}) {
	logCore(l, lvl, "", args, skipFramesNum)
}

func (l *Logger) Errorv(args ...interface{}) {
	logCore(l, ErrorLevel, "", args, skipFramesNum)
}

func (l *Logger) Warnv(args ...interface{}) {
	logCore(l, WarnLevel, "", args, skipFramesNum)
}

func (l *Logger) Infov(args ...interface{}) {
	logCore(l, InfoLevel, "", args, skipFramesNum)
}

func (l *Logger) Debugv(args ...interface{}) {
	logCore(l, DebugLevel, "", args, skipFramesNum)
}

func (l *Logger) Tracev(args ...interface{}) {
	logCore(l, TraceLevel, "", args, skipFramesNum)
}

// Logf methods

func (l *Logger) Logf(lvl Level, format string, args ...interface{}) {
	logCore(l, lvl, format, args, skipFramesNum)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	logCore(l, ErrorLevel, format, args, skipFramesNum)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	logCore(l, WarnLevel, format, args, skipFramesNum)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	logCore(l, InfoLevel, format, args, skipFramesNum)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	logCore(l, DebugLevel, format, args, skipFramesNum)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	logCore(l, TraceLevel, format, args, skipFramesNum)
}

func (l *Logger) IsEnabled(lvl Level) bool {
	return l.filter(lvl)
}

func (l *Logger) Flush() {
	l.sink.Flush()
}

func (l *Logger) With(name string, value interface{}) *Logger {
	s := make([]Field, 0, len(l.fields)+1)
	s = append(s, l.fields...)
	s = append(s, Field{name, value})
	return &Logger{l.sink, l.filter, s, l.flags}
}

func (l *Logger) WithFields(fields []Field) *Logger {
	s := make([]Field, 0, len(l.fields)+len(fields))
	s = append(s, l.fields...)
	s = append(s, fields...)
	return &Logger{l.sink, l.filter, s, l.flags}
}
