package log

import "io"

type Recorder interface {
	Log(*Record)
	Flush()
	IsEnabled(*Metadata) bool
}

type NopRecorder struct{}

func (NopRecorder) Log(*Record) {
}

func (NopRecorder) Flush() {
}
func (NopRecorder) IsEnabled(*Metadata) bool {
	return false
}

type LeveledRecorder struct {
	MaxLevel Level
	Target   Recorder
}

func (l *LeveledRecorder) Log(r *Record) {
	l.Target.Log(r)
}

func (l *LeveledRecorder) Flush() {
}

func (l *LeveledRecorder) IsEnabled(m *Metadata) bool {
	return m.Level <= l.MaxLevel
}

type StreamRecorder struct {
	Formatter func(*Record) string
	Stream    io.Writer
}

func (s *StreamRecorder) Log(r *Record) {
	io.WriteString(s.Stream, s.Formatter(r))
}

func (s *StreamRecorder) Flush() {
}

func (s *StreamRecorder) IsEnabled(m *Metadata) bool {
	return true
}

type MultiRecorder []Recorder

func (m MultiRecorder) Log(r *Record) {
	for _, c := range m {
		if c.IsEnabled(&r.Meta) {
			c.Log(r)
		}
	}
}

func (m MultiRecorder) Flush() {
	for _, c := range m {
		c.Flush()
	}
}

func (m MultiRecorder) IsEnabled(meta *Metadata) bool {
	for _, c := range m {
		if c.IsEnabled(meta) {
			return true
		}
	}
	return false
}

func CreateMultiRecorder(recs ...Recorder) Recorder {
	return MultiRecorder(recs)
}
