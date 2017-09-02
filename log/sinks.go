package log

import "io"
import "sync"

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
	Inner    Sink
}

func (l *LeveledSink) Log(r *Record) {
	l.Inner.Log(r)
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

func (s *SyncedSink) IsEnabled(m *Metadata) bool {
	return s.Inner.IsEnabled(m)
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
