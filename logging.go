package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func CreateLogStream(logFile *string, overwriteLog bool) io.WriteCloser {
	if logFile == nil || *logFile == "" {
		return os.Stderr
	}
	fileName, err := filepath.Abs(*logFile)
	if err != nil {
		log.Println(err.Error())
		return os.Stderr
	}
	var flag int
	if overwriteLog {
		flag = os.O_CREATE | os.O_TRUNC
	} else {
		flag = os.O_CREATE | os.O_APPEND
	}
	file, err := os.OpenFile(fileName, flag, 0)
	if err != nil {
		log.Println(err.Error())
		return os.Stderr
	}
	return file
}
