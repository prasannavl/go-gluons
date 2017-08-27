package main

import (
	"context"
	"net/http"
	"pvl/apicore/appcontext"
	"time"

	flag "github.com/spf13/pflag"

	"fmt"

	"pvl/apicore/platform"

	"github.com/prasannavl/go-grab/lifecycle"
	"github.com/prasannavl/go-grab/log"
	logc "github.com/prasannavl/go-grab/log-config"
)

func main() {
	var addr string
	var logFile string
	var logDisabled bool
	var verbosity int

	platform.Init()

	flag.Usage = func() {
		fmt.Printf("\nUsage: [opts]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.CountVarP(&verbosity, "verbose", "v", "verbosity level")
	flag.StringVarP(&addr, "address", "a", "localhost:8000", "the 'host:port' for the service to listen on")
	flag.StringVar(&logFile, "log", "", "the log file destination")
	flag.BoolVar(&logDisabled, "no-log", false, "disable the logger")
	flag.Parse()

	logInitResult := logc.LogInitResult{}
	if !logDisabled {
		logOpts := logc.DefaultOptions()
		if logFile != "" {
			logOpts.LogFile = logFile
		}
		logOpts.VerbosityLevel = verbosity
		logc.Init(&logOpts, &logInitResult)
	}

	log.Infof("listen-address: %q", addr)

	c := createAppContext(logInitResult.Logger, addr)
	runServer(c, NewApp(c))
}

func runServer(c *appcontext.AppContext, handler http.Handler) {
	server := &http.Server{
		Addr:           c.ServerAddress,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	lifecycle.CreateShutdownHandler(func() {
		server.Shutdown(context.Background())
	}, lifecycle.ShutdownSignals...)

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Errorf("server close: %s", err.Error())
	}
	log.Info("exit")
}
