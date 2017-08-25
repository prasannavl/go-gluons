package platform

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/prasannavl/go-grab/log"
)

func SetExitHandler(handler func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

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
