package log_test

import (
	"os"
	"testing"

	"github.com/prasannavl/go-grab/log"
)

func TestPrint(t *testing.T) {
	rec := log.CreateMultiRecorder(
		&log.LeveledRecorder{
			MaxLevel: log.InfoLevel,
			Target: &log.StreamRecorder{
				Formatter: log.DefaultColorTextFormatterForHuman,
				Stream:    os.Stdout,
			},
		},
	)

	l := log.New(rec)
	log.SetGlobal(l)

	log.Info("Hello there 1")
	log.Warn("Hello there 2")
	log.Error("Hello there 3")
	log.Debug("Hello there 4")
	log.Trace("Hello there 5")

	log.Infof("%s", "Hey you X")
	log.Warnf("%s %q %v", "Hey", "you", "Y")

	l2 := log.WithContext("ctxName", "some val")
	l2.Info("hello there!!")
	l2.WithContext("ctx2", "another val").Info("Hey you")
	l2.Infof("%s %v", "hello there", "again")
}
