package log_test

import (
	"os"
	"testing"

	"github.com/prasannavl/go-gb/log"
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

	log.SetGlobal(rec)

	log.Info("Hello there 1")
	log.Warn("Hello there 2")
	log.Error("Hello there 3")
	log.Debug("Hello there 4")
	log.Trace("Hello there 5")

	log.Infof("%s", "Hey you X")
	log.Warnf("%s %q %v", "Hey", "you", "Y")

	l := log.WithContext("ctxName", "some val")
	l.Info("hello there!!")
	l.Infof("%s %v", "hello there", "again")
}
