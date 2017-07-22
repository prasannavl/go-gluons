package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"encoding/json"
	"fmt"
	"time"
)

func createLogStream(logFile *string, overwriteLog bool) io.WriteCloser {
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

//go:generate echo Hello, Go Generate!
func main() {
	var addr string
	var logFile string
	var logOverwrite bool
	var logDisable bool

	flag.Usage = func() {
		fmt.Printf("\nUsage: [opts]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.StringVar(&addr, "address", "localhost:8000", "the `address:port` for the service to listen on")
	flag.StringVar(&logFile, "log", "", "the log file destination")
	flag.BoolVar(&logOverwrite, "log-overwrite", false, "overwrite the log file if set, or appends")
	flag.BoolVar(&logDisable, "log-disable", false, "disable the logger")
	flag.Parse()

	if logDisable {
		log.SetOutput(ioutil.Discard)
	} else {
		logWriter := createLogStream(&logFile, logOverwrite)
		defer logWriter.Close()
		log.SetOutput(logWriter)
	}

	log.Printf("Listen address: %s", addr)

	r := func(w http.ResponseWriter, req *http.Request) {
		data := struct {
			Name string
			Date time.Time
		}{"Hello", time.Now()}
		j, _ := json.Marshal(&data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}

	log.Println("Listening..")
	http.HandleFunc("/", r)

	server := &http.Server{Addr: addr, Handler: http.DefaultServeMux}
	var _ = server.ListenAndServe()
	log.Println("Exit")
	time.Sleep(3 * time.Second)
}
