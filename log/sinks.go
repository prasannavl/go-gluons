package log

import "io"
import "sync"

type Sink interface {
	Log(*Record)
	Flush()
}

type NopSink struct{}

func (NopSink) Log(*Record) {}

func (NopSink) Flush() {}

type LeveledSink struct {
	MaxLevel Level
	Inner    Sink
}

func (l *LeveledSink) Log(r *Record) {
	if r.Meta.Level <= l.MaxLevel {
		l.Inner.Log(r)
	}
}

func (l *LeveledSink) Flush() {
	l.Inner.Flush()
}

type StreamSink struct {
	Formatter func(*Record) string
	Stream    io.Writer
}

func (s *StreamSink) Log(r *Record) {
	io.WriteString(s.Stream, s.Formatter(r))
}

func (s *StreamSink) Flush() {}

type SyncedSink struct {
	m     sync.Mutex
	Inner Sink
}

func (s *SyncedSink) Log(r *Record) {
	s.m.Lock()
	defer s.m.Unlock()
	s.Inner.Log(r)
}

func (s *SyncedSink) Flush() {
	s.m.Lock()
	defer s.m.Unlock()
	s.Inner.Flush()
}

type MultiSink []Sink

func (m MultiSink) Log(r *Record) {
	for _, c := range m {
		c.Log(r)
	}
}

func (m MultiSink) Flush() {
	for _, c := range m {
		c.Flush()
	}
}

func CreateMultiSink(recs ...Sink) Sink {
	return MultiSink(recs)
}
