package log

type LogWriter struct {
	l          *Logger
	WriteLevel Level
	Prefix     string
}

func (s *LogWriter) Write(p []byte) (n int, err error) {
	var m string
	if len(s.Prefix) > 0 {
		m = s.Prefix + string(p)
	} else {
		m = string(p)
	}
	s.l.Log(s.WriteLevel, m)
	return len(p), nil
}

func NewLogWriter(l *Logger, lvl Level, prefix string) *LogWriter {
	return &LogWriter{l, lvl, prefix}
}
