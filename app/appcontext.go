package app

import (
	"github.com/prasannavl/go-grab/log"
)

type AppContext struct {
	Services
	ServerAddress string
}

// Services is the global services context
type Services struct {
	Logger *log.Logger
}
