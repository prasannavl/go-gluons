package lifecycle

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/prasannavl/go-grab/log"
)

var ShutdownSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP}

func CreateShutdownHandler(handler func(), s ...os.Signal) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, s...)

	go func() {
		for sig := range quit {
			log.Infof("signal: %v", sig)
			log.Info("shutting down..")
			if handler != nil {
				handler()
			}
		}
	}()
}
