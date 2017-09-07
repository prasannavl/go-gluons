package log_test

import (
	"os"
	"testing"

	"github.com/prasannavl/go-gluons/log"
)

func TestPrint(t *testing.T) {

	sink := log.CreateMultiSink(
		&log.LeveledSink{
			MaxLevel: log.InfoLevel,
			Inner: &log.StreamSink{
				Formatter: log.DefaultTextFormatterForHuman,
				Stream:    os.Stdout,
			},
		},
		&log.LeveledSink{
			MaxLevel: log.ErrorLevel,
			Inner: &log.StreamSink{
				Formatter: log.DefaultTextFormatter,
				Stream:    os.Stdout,
			},
		},
	)

	l := log.New(sink)
	log.SetLogger(l)

	log.Info("Hello there 1")
	log.Warn("Hello there 2")
	log.Error("Hello there 3")
	log.Debug("Hello there 4")
	log.Trace("Hello there 5")

	log.Infof("%s", "Hey you X")
	log.Warnf("%s %q %v", "Hey", "you", "Y")

	l2 := log.With("ctxName", "some val")
	l2.Info("hello there!!")
	l2.With("ctx2", "another val").Info("Hey you")
	l2.Infof("%s %v", "hello there", "again")
}
