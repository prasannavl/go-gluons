package appcontext

import "go.uber.org/zap"

type AppContext struct {
	Services
}

// Services is the global services context
type Services struct {
	Logger *zap.Logger
}
