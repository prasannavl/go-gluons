package log

import "io"

type Sink interface {
	Log(*Record)
	Flush()
	IsEnabled(*Metadata) bool
}

type NopSink struct{}

func (NopSink) Log(*Record) {
}

func (NopSink) Flush() {
}
func (NopSink) IsEnabled(*Metadata) bool {
	return false
}

type LeveledSink struct {
	MaxLevel Level
	Target   Sink
}

func (l *LeveledSink) Log(r *Record) {
	l.Target.Log(r)
}

func (l *LeveledSink) Flush() {
}

func (l *LeveledSink) IsEnabled(m *Metadata) bool {
	return m.Level <= l.MaxLevel
}

type StreamSink struct {
	Formatter func(*Record) string
	Stream    io.Writer
}

func (s *StreamSink) Log(r *Record) {
	io.WriteString(s.Stream, s.Formatter(r))
}

func (s *StreamSink) Flush() {
}

func (s *StreamSink) IsEnabled(m *Metadata) bool {
	return true
}

type MultiSink []Sink

func (m MultiSink) Log(r *Record) {
	for _, c := range m {
		if c.IsEnabled(&r.Meta) {
			c.Log(r)
		}
	}
}

func (m MultiSink) Flush() {
	for _, c := range m {
		c.Flush()
	}
}

func (m MultiSink) IsEnabled(meta *Metadata) bool {
	for _, c := range m {
		if c.IsEnabled(meta) {
			return true
		}
	}
	return false
}

func CreateMultiSink(recs ...Sink) Sink {
	return MultiSink(recs)
}
