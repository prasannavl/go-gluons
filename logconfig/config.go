package logconfig

import (
	"io"
	stdlog "log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/prasannavl/go-grab/log"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	VerbosityLevel   int
	LogFile          string
	FallbackFileName string
	FallbackDir      string
	LoggerMutex      bool

	Rolling         bool
	MaxSize         int // megabytes
	MaxBackups      int
	MaxAge          int // days
	CompressBackups bool
	Humanize        int
	StdLogLevel     log.Level
}

func DefaultOptions() Options {
	return Options{
		VerbosityLevel:   VerbosityLevel.Warn,
		LogFile:          CommonTargets.TargetStdOut,
		FallbackFileName: "run.log",
		FallbackDir:      "logs",
		Rolling:          true,
		LoggerMutex:      false,
		MaxSize:          100,
		MaxBackups:       2,
		MaxAge:           28,
		CompressBackups:  true,
		Humanize:         Humanize.Auto,
		StdLogLevel:      log.TraceLevel,
	}
}

func Init(opts *Options, result *LogInitResult) {
	result.Enabled = false
	logFile := opts.LogFile
	if logFile == CommonTargets.TargetNull {
		return
	}
	level := logLevelFromVerbosityLevel(opts.VerbosityLevel)
	if level == 0 {
		return
	}
	s, name := mustCreateWriteStream(opts)
	var formatter func(r *log.Record) string

	if opts.Humanize == Humanize.True {
		formatter = log.DefaultColorTextFormatterForHuman
	} else {
		formatter = log.DefaultTextFormatter
	}

	var target log.Sink

	target = &log.StreamSink{
		Formatter: formatter,
		Stream:    s,
	}

	if opts.LoggerMutex {
		target = &log.SyncedSink{
			Inner: target,
		}
	}

	target = &log.LeveledSink{
		MaxLevel: level,
		Target:   target,
	}

	l := log.New(target)
	log.SetGlobal(l)
	stdWriter := log.NewLogWriter(l, opts.StdLogLevel, "std: ")
	stdlog.SetOutput(stdWriter)

	result.Enabled = true
	result.Filename = name
	result.Logger = l
	result.Writer = s
	result.StdWriter = stdWriter
	result.StdLogger = stdlog.New(stdWriter, "", 0)
}

type LogInitResult struct {
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

func mustCreateWriteStream(opts *Options) (w io.Writer, filename string) {
	var err error
	logFile := opts.LogFile
	const errFormat = "error: logger => %s"
	if logFile == "" {
		logFile, err = checkedLogFileName(filepath.Clean(opts.FallbackDir + "/" + opts.FallbackFileName))
		if err != nil {
			stdlog.Fatalf(errFormat, err.Error())
		}
	}
	switch logFile {
	case CommonTargets.TargetStdOut:
		return os.Stdout, logFile
	case CommonTargets.TargetStdErr:
		return os.Stderr, logFile
	default:
		if err := ensureFileParentDir(logFile); err != nil {
			stdlog.Fatalf(errFormat, err.Error())
		}
		if logFile, err = checkedLogFileName(logFile); err != nil {
			stdlog.Fatalf(errFormat, err.Error())
		}
		if !opts.Rolling {
			fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_EXCL, os.FileMode(0644))
			if err != nil {
				stdlog.Println(errFormat, err.Error())
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

// This method tries to touch the file, and if not,
// one possibility is that it's being used by another
// process. So, try again once with the
// PID appended. If that fails too - error.
func checkedLogFileName(logFile string) (string, error) {
	filename := logFile
	if err := touchFile(filename); err != nil {
		filename = filename + ".pid-" + strconv.Itoa(os.Getpid()) + ".txt"
		if e := touchFile(filename); e != nil {
			// Return the old error
			return "", err
		}
	}
	return filename, nil
}

func ensureFileParentDir(path string) error {
	d := filepath.Dir(path)
	err := os.MkdirAll(d, os.FileMode(0777))
	if err != nil {
		return err
	}
	return nil
}

func touchFile(path string) error {
	var err error
	if err = ensureFileParentDir(path); err != nil {
		return err
	}
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_EXCL, os.FileMode(0644))
	if err != nil {
		return err
	}
	if err = fd.Close(); err != nil {
		return err
	}
	return nil
}

// Enums

type (
	humanizeEnum struct {
		Auto  int
		False int
		True  int
	}

	commonTargetEnum struct {
		TargetStdOut string
		TargetStdErr string
		TargetNull   string
	}

	verbosityLevel struct {
		Error int
		Warn  int
		Info  int
		Debug int
		Trace int
	}
)

var (
	Humanize = humanizeEnum{
		Auto:  -1,
		False: 0,
		True:  1,
	}

	CommonTargets = commonTargetEnum{
		TargetStdOut: ":stdout",
		TargetStdErr: ":stderr",
		TargetNull:   ":null",
	}

	VerbosityLevel = verbosityLevel{
		Error: -1,
		Warn:  0,
		Info:  1,
		Debug: 2,
		Trace: 3,
	}
)
