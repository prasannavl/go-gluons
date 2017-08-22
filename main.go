package main

import (
	"context"
	"net/http"
	"os"
	"pvl/apicore/appcontext"
	"pvl/apicore/logger"
	"pvl/apicore/platform"
	"time"

	flag "github.com/spf13/pflag"

	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	var addr string
	var logFile string
	var logDisabled bool
	var debugMode bool

	platform.Init()

	flag.Usage = func() {
		fmt.Printf("\nUsage: [opts]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.StringVarP(&addr, "address", "a", "localhost:8000", "the 'host:port' for the service to listen on")
	flag.StringVar(&logFile, "log-file", "", "the log file destination")
	flag.BoolVar(&logDisabled, "log-off", false, "disable the logger")
	flag.Parse()

	log := logger.Create(!logDisabled, logFile, debugMode)
	zap.RedirectStdLog(log)

	log.Info("args", zap.String("listen-address", addr))

	c := createAppContext(log, addr)
	runServer(c, NewApp(c))
}

func runServer(c *appcontext.AppContext, handler http.Handler) {
	log := c.Services.Logger
	server := &http.Server{
		Addr:           c.ServerAddress,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range quit {
			log.Info(fmt.Sprintf("signal: %v", sig))
			log.Info("shutting down..")
			server.Shutdown(context.Background())
		}
	}()

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal("server close", zap.Error(err))
	}
	log.Info("exit")
}
