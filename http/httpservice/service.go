package httpservice

import (
	"context"
	"net"
	"net/http"
	"time"
)

type HttpService struct {
	name      string
	server    *http.Server
	l         net.Listener
	isRunning bool
}

func New(name string, server *http.Server, listener net.Listener) Service {
	return &HttpService{name: name, server: server, l: listener}
}

func (r *HttpService) IsRunning() bool {
	return r.isRunning
}

func (r *HttpService) Start() {
	if !r.isRunning {
		r.isRunning = true
		r.server.Serve(r.l)
	}
}

func (r *HttpService) Name() string {
	return r.name
}

func (r *HttpService) Stop(timeout time.Duration) {
	if r.isRunning {
		var ctx context.Context
		var cancel context.CancelFunc
		if timeout == 0 {
			ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(timeout))
		} else {
			ctx = context.Background()
		}
		r.server.Shutdown(ctx)
		if cancel != nil {
			cancel()
		}
	}
}

type Service interface {
	Name() string
	IsRunning() bool
	Start()
	Stop(timeout time.Duration)
}
