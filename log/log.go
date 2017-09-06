package log

import (
	"os"
	"time"
)

type Level uint

const (
	_                = 0
	ErrorLevel Level = 1 << (iota - 1)
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger struct {
	sink   Sink
	fields []Field
}

type Record struct {
	l      Logger
	Meta   Metadata
	Format string
	Args   []interface{}
}

func (r *Record) Fields() []Field {
	return r.l.fields
}

type Field struct {
	Name  string
	Value interface{}
}

type Metadata struct {
	Level Level
	Time  time.Time
	File  string
	Line  uint
}

func newRecord(l Logger, m Metadata, format string, args []interface{}) Record {
	return Record{l: l, Meta: m, Format: format, Args: args}
}

func newMetadata(lvl Level, withTime bool) Metadata {
	if withTime {
		return Metadata{Level: lvl, Time: time.Now()}
	}
	return Metadata{Level: lvl}
}

func (l Logger) IsEnabled(lvl Level) bool {
	m := newMetadata(lvl, false)
	return l.sink.IsEnabled(&m)
}

func (l Logger) Logf(lvl Level, format string, args ...interface{}) {
	if !l.IsEnabled(lvl) {
		return
	}
	r := newRecord(l, newMetadata(lvl, true), format, args)
	l.sink.Log(&r)
}

func (l Logger) Log(lvl Level, message string) {
	l.Logf(lvl, message)
}

func (l Logger) Logv(lvl Level, args ...interface{}) {
	l.Logf(lvl, "", args...)
}

func (l Logger) Info(message string) {
	l.Log(InfoLevel, message)
}

func (l Logger) Warn(message string) {
	l.Log(WarnLevel, message)
}

func (l Logger) Error(message string) {
	l.Log(ErrorLevel, message)
}

func (l Logger) Debug(message string) {
	l.Log(DebugLevel, message)
}

func (l Logger) Trace(message string) {
	l.Log(TraceLevel, message)
}

func (l Logger) Infov(args ...interface{}) {
	l.Logv(InfoLevel, args...)
}

func (l Logger) WarnV(args ...interface{}) {
	l.Logv(WarnLevel, args...)
}

func (l Logger) Errorv(args ...interface{}) {
	l.Logv(ErrorLevel, args...)
}

func (l Logger) Debugv(args ...interface{}) {
	l.Logv(DebugLevel, args...)
}

func (l Logger) Tracev(args ...interface{}) {
	l.Logv(TraceLevel, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.Logf(InfoLevel, format, args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	l.Logf(WarnLevel, format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Logf(ErrorLevel, format, args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.Logf(DebugLevel, format, args...)
}

func (l Logger) Tracef(format string, args ...interface{}) {
	l.Logf(TraceLevel, format, args...)
}

func (l Logger) Flush() {
	l.sink.Flush()
}

func (l Logger) With(name string, value interface{}) Logger {
	s := make([]Field, 0, len(l.fields)+1)
	s = append(s, l.fields...)
	s = append(s, Field{name, value})
	return Logger{l.sink, s}
}

func (l Logger) WithFields(fields []Field) Logger {
	s := make([]Field, 0, len(l.fields)+len(fields))
	s = append(s, l.fields...)
	s = append(s, fields...)
	return Logger{l.sink, s}
}

func IsEnabled(lvl Level) bool {
	return g.IsEnabled(lvl)
}

func Logf(lvl Level, format string, args ...interface{}) {
	g.Logf(lvl, format, args...)
}

func Log(lvl Level, message string) {
	g.Log(lvl, message)
}

func Logv(lvl Level, args ...interface{}) {
	g.Logv(lvl, args...)
}

func Info(message string) {
	g.Info(message)
}

func Warn(message string) {
	g.Warn(message)
}

func Error(message string) {
	g.Error(message)
}

func Debug(message string) {
	g.Debug(message)
}

func Trace(message string) {
	g.Trace(message)
}

func Infov(args ...interface{}) {
	g.Infov(args...)
}

func WarnV(args ...interface{}) {
	g.WarnV(args...)
}

func Errorv(args ...interface{}) {
	g.Errorv(args...)
}

func Debugv(args ...interface{}) {
	g.Debugv(args...)
}

func Tracev(args ...interface{}) {
	g.Tracev(args...)
}

func Infof(format string, args ...interface{}) {
	g.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	g.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	g.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	g.Debugf(format, args...)
}

func Tracef(format string, args ...interface{}) {
	g.Tracef(format, args...)
}

func Flush() {
	g.Flush()
}

func With(name string, val interface{}) Logger {
	return g.With(name, val)
}

func WithFields(fields []Field) Logger {
	return g.WithFields(fields)
}

var (
	g *Logger
)

func init() {
	g = &Logger{sink: &LeveledSink{
		MaxLevel: ErrorLevel,
		Inner: &StreamSink {
			Stream: os.Stderr,
			Formatter: DefaultColorTextFormatterForHuman,
		},
	}}
}

func SetLogger(l *Logger) {
	g = l
}

func GetLogger() *Logger {
	return g
}

func New(sink Sink) *Logger {
	return &Logger{sink, nil}
}

func GetSink(logger *Logger) Sink {
	return logger.sink
}
