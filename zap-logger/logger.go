package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func Create(isEnabled bool, logFile string, debugMode bool, fallBackFileName string) *zap.Logger {
	if !isEnabled || logFile == ":null" {
		return createDisabledLogger(debugMode)
	}
	core := createZapProductionCore(createWriteStream(logFile, fallBackFileName), debugMode)
	if debugMode {
		dcore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(newEncoderConfig(configFlagDev|configFlagConsole)),
			os.Stderr, zap.DebugLevel)
		return zap.New(zapcore.NewTee(dcore, core))
	}

	core = zapcore.NewSampler(core, time.Second, 100, 100)
	return zap.New(core)
}

func createWriteStream(logFile string, fallBackFileName string) io.Writer {
	if logFile == "" {
		logDir := "logs/"
		defaultLog := logDir + fallBackFileName
		_ = os.MkdirAll(logDir, os.FileMode(777))
		f, err := os.OpenFile(defaultLog, os.O_CREATE, os.FileMode(0644))
		if err != nil {
			return os.Stderr
		}
		f.Close()
		logFile = defaultLog
	}
	switch logFile {
	case ":stdout":
		return os.Stdout
	case ":stderr":
		return os.Stderr
	default:
		fileName, err := filepath.Abs(logFile)
		if err != nil {
			return os.Stderr
		}
		return &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}
	}
}

func createDisabledLogger(debugMode bool) *zap.Logger {
	if debugMode {
		return mustCreateDevLogger()
	}
	return zap.NewNop()
}

func createZapProductionCore(writer io.Writer, debugMode bool) zapcore.Core {
	var level zapcore.LevelEnabler
	if debugMode {
		level = zap.DebugLevel
	} else {
		level = zap.InfoLevel
	}
	w := zapcore.AddSync(writer)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig(configFlagNone)),
		w, level)

	return core
}

func mustCreateDevLogger() *zap.Logger {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig = newEncoderConfig(configFlagDev | configFlagConsole)
	l, err := c.Build()
	if err != nil {
		panic(err)
	}
	return l
}

type configFlags uint

const (
	configFlagNone configFlags = 0
	configFlagDev              = 1 << iota
	configFlagConsole
)

func newEncoderConfig(flags configFlags) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if flags|configFlagDev == flags {
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncodeDuration = zapcore.StringDurationEncoder
	}
	if flags|configFlagConsole == flags {
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}
	return config
}
