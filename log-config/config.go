package logger

import (
	"io"
	stdlog "log"
	"os"
	"path/filepath"

	"github.com/prasannavl/go-gb/log"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	VerbosityLevel   int
	LogFile          string
	FallbackFileName string
	FallbackDir      string

	Rolling         bool
	MaxSize         int // megabytes
	MaxBackups      int
	MaxAge          int // days
	CompressBackups bool
	NoColor         bool
	Humanize        bool
}

func DefaultOptions() Options {
	return Options{
		VerbosityLevel:   0,
		LogFile:          "",
		FallbackFileName: "run.log",
		FallbackDir:      "logs",
		Rolling:          true,
		MaxSize:          100,
		MaxBackups:       2,
		MaxAge:           28,
		CompressBackups:  true,
		NoColor:          false,
		Humanize:         true,
	}
}

const (
	StdOut   = ":stdout"
	StdErr   = ":stderr"
	Disabled = ":null"
)

func Init(opts *Options) {
	logFile := opts.LogFile
	if logFile == Disabled {
		return
	}
	level := logLevelFromVerbosityLevel(opts.VerbosityLevel)
	if level == 0 {
		return
	}
	s := createWriteStream(opts)
	var formatter func(r *log.Record) string

	if opts.Humanize {
		if opts.NoColor {
			formatter = log.DefaultTextFormatterForHuman
		} else {
			formatter = log.DefaultColorTextFormatterForHuman
		}
	} else {
		formatter = log.DefaultTextFormatter
	}

	target := log.StreamRecorder{
		Formatter: formatter,
		Stream:    s,
	}

	rec := log.LeveledRecorder{
		MaxLevel: log.InfoLevel,
		Target:   &target,
	}

	log.SetGlobal(&rec)
}

func logLevelFromVerbosityLevel(vLevel int) log.Level {
	switch vLevel {
	case -1:
		return log.ErrorLevel
	case 0:
		return log.WarnLevel
	case 1:
		return log.InfoLevel
	case 2:
		return log.DebugLevel
	case 3:
		return log.TraceLevel
	}
	return log.TraceLevel
}

func createWriteStream(opts *Options) io.Writer {
	var err error
	logFile := opts.LogFile
	const loggerErrFormat = "error: logger => %s"
	if logFile == "" {
		if logFile, err = touchFile(opts.FallbackDir, opts.FallbackFileName); err != nil {
			stdlog.Fatalf(loggerErrFormat, err.Error())
		}
	}
	switch logFile {
	case StdOut:
		return os.Stdout
	case StdErr:
		return os.Stderr
	default:
		if err := touchFilePath(logFile); err != nil {
			stdlog.Fatalf(loggerErrFormat, err.Error())
		}
		if !opts.Rolling {
			fd, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE, os.FileMode(0644))
			if err != nil {
				stdlog.Fatalf(loggerErrFormat, err.Error())
			}
			return fd
		}
		return &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.CompressBackups,
		}
	}
}

func touchFilePath(path string) error {
	var err error
	a, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	d := filepath.Dir(a)
	err = os.MkdirAll(d, os.FileMode(0777))
	if err != nil {
		return err
	}
	return nil
}

func touchFile(dir string, filename string) (string, error) {
	var err error
	d := filepath.Clean(dir)
	f := d + "/" + filepath.Clean(filename)
	err = os.MkdirAll(d, os.FileMode(0777))
	if err != nil {
		return "", err
	}
	fd, err := os.OpenFile(f, os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return "", err
	}
	if err = fd.Close(); err != nil {
		return "", err
	}
	return f, nil
}
