package logger

import (
	"io"
	stdlog "log"
	"os"
	"path/filepath"

	"github.com/prasannavl/go-grab/log"
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
	StdLogLevel     log.Level
}

func DefaultOptions() Options {
	return Options{
		VerbosityLevel:   0,
		LogFile:          TargetStdOut,
		FallbackFileName: "run.log",
		FallbackDir:      "logs",
		Rolling:          true,
		MaxSize:          100,
		MaxBackups:       2,
		MaxAge:           28,
		CompressBackups:  true,
		NoColor:          false,
		Humanize:         true,
		StdLogLevel:      log.TraceLevel,
	}
}

const (
	TargetStdOut = ":stdout"
	TargetStdErr = ":stderr"
	TargetNull   = ":null"
)

func Init(opts *Options, meta *LogInstanceMeta) {
	meta.Enabled = false
	logFile := opts.LogFile
	if logFile == TargetNull {
		return
	}
	level := logLevelFromVerbosityLevel(opts.VerbosityLevel)
	if level == 0 {
		return
	}
	s, name := mustCreateWriteStream(opts)
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
		MaxLevel: level,
		Target:   &target,
	}

	l := log.New(&rec)
	log.SetGlobal(l)
	stdWriter := log.NewLogWriter(l, opts.StdLogLevel, "std: ")
	stdlog.SetOutput(stdWriter)

	meta.Enabled = true
	meta.Filename = name
	meta.Logger = l
	meta.Writer = s
	meta.StdWriter = stdWriter
	meta.StdLogger = stdlog.New(stdWriter, "", 0)
}

type LogInstanceMeta struct {
	Enabled   bool
	Filename  string
	Writer    io.Writer
	Logger    *log.Logger
	StdWriter *log.LogWriter
	StdLogger *stdlog.Logger
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

// TODO: Handle the case of parallel logs using exclusive write streams, and
// create logs suffixed with the datetime to create unique paths,
// before lumberjack fails
func mustCreateWriteStream(opts *Options) (w io.Writer, filename string) {
	var err error
	logFile := opts.LogFile
	const errFormat = "error: logger => %s"
	if logFile == "" {
		if logFile, err = touchFile(opts.FallbackDir, opts.FallbackFileName); err != nil {
			stdlog.Fatalf(errFormat, err.Error())
		}
	}
	switch logFile {
	case TargetStdOut:
		return os.Stdout, TargetStdOut
	case TargetStdErr:
		return os.Stderr, TargetStdErr
	default:
		if err := touchFilePath(logFile); err != nil {
			stdlog.Fatalf(errFormat, err.Error())
		}
		if !opts.Rolling {
			fd, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE, os.FileMode(0644))
			if err != nil {
				stdlog.Fatalf(errFormat, err.Error())
			}
			return fd, logFile
		}
		return &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.CompressBackups,
		}, logFile
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
