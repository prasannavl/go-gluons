package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"fmt"
	"os/signal"
	"syscall"
)

func main() {
	var addr string
	var logFile string
	var logOverwrite bool
	var logDisable bool
	var debugMode bool

	flag.Usage = func() {
		fmt.Printf("\nUsage: [opts]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.StringVar(&addr, "address", "localhost:8000", "the `address:port` for the service to listen on")
	flag.StringVar(&logFile, "log", "", "the log file destination")
	flag.BoolVar(&logOverwrite, "log-overwrite", false, "overwrite the log file if set, or appends")
	flag.BoolVar(&logDisable, "log-disable", false, "disable the logger")
	flag.Parse()

	logStream := CreateLogStream(&logFile, logOverwrite, logDisable, debugMode)
	defer logStream.Close()
	log.SetOutput(logStream)

	log.Printf("Listen address: %s", addr)

	runServer(addr, handler(addr))
}

func handler(addr string) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("labs.prasannavl.com/", NewHandlerFunc("PVL Labs", nil))
	m.HandleFunc(addr+"/", NewHandlerFunc(addr, nil))
	m.HandleFunc("nf.statwick.com/", NewHandlerFunc("NextFirst API", nil))
	return m
}

func runServer(addr string, handler http.Handler) {
	server := &http.Server{Addr: addr, Handler: handler}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range quit {
			log.Printf("Signal: %q", sig)
			log.Printf("Shutting down..")
			server.Shutdown(context.Background())
		}
	}()

	log.Println("Listening..")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Fatal: %v", err)
	}
	log.Println("Exit: Success")
}
