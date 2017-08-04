package main

import (
	"context"
	"log"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"

	"fmt"
	"os/signal"
	"syscall"
)

func main() {
	var addr string
	var logFile string
	var logOverwrite bool
	var logDisabled bool
	var debugMode bool

	flag.Usage = func() {
		fmt.Printf("\nUsage: [opts]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.StringVarP(&addr, "address", "a", "localhost:8000", "the 'host:port' for the service to listen on")
	flag.StringVar(&logFile, "log-file", "", "the log file destination")
	flag.BoolVar(&logOverwrite, "log-overwrite", false, "overwrite the log file if set, or appends")
	flag.BoolVar(&logDisabled, "log-off", false, "disable the logger")
	flag.Parse()

	logStream := CreateLogStream(&logFile, logOverwrite, logDisabled, debugMode)
	defer logStream.Close()
	log.SetOutput(logStream)

	log.Printf("listen address: %s", addr)

	runServer(addr, NewApp(addr))
}

func runServer(addr string, handler http.Handler) {
	server := &http.Server{Addr: addr, Handler: handler}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range quit {
			log.Printf("signal: %q", sig)
			log.Printf("shutting down..")
			server.Shutdown(context.Background())
		}
	}()

	log.Println("listening..")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("fatal: %v", err)
	}
	log.Println("exit: ok")
}
