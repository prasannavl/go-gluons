package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

type writeCloser struct {
	io.Writer
	io.Closer
}

func CreateLogStream(logFile *string, truncate bool, disabled bool, debugMode bool) io.WriteCloser {
	w, closer := createLogStream(logFile, truncate, disabled, debugMode)
	if closer == nil {
		return nopWriteCloser{w}
	}
	return writeCloser{w, closer}
}

func createLogStream(logFile *string, truncate bool, disabled bool, debugMode bool) (io.Writer, io.Closer) {
	if disabled {
		if debugMode {
			return os.Stderr, nil
		}
		return ioutil.Discard, nil
	}

	if logFile == nil || *logFile == "" {
		return os.Stderr, nil
	}

	fileName, err := filepath.Abs(*logFile)
	if err != nil {
		log.Println(err.Error())
		return os.Stderr, nil
	}
	var flag int
	if truncate {
		flag = os.O_CREATE | os.O_TRUNC
	} else {
		flag = os.O_CREATE | os.O_APPEND
	}
	file, err := os.OpenFile(fileName, flag, 0)
	if err != nil {
		log.Println(err.Error())
		return os.Stderr, nil
	}

	if debugMode {
		return io.MultiWriter(file, os.Stderr), file
	}
	return file, file
}
