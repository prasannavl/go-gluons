package httpservice

import (
	"context"
	"errors"
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

func (r *HttpService) Run() error {
	if !r.isRunning {
		r.isRunning = true
		return r.server.Serve(r.l)
	}
	return errors.New("service attempted to start while already running")
}

func (r *HttpService) Name() string {
	return r.name
}

func (r *HttpService) Stop(timeout time.Duration) error {
	if r.isRunning {
		var ctx context.Context
		var cancel context.CancelFunc
		if timeout == 0 {
			ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(timeout))
		} else {
			ctx = context.Background()
		}
		err := r.server.Shutdown(ctx)
		if cancel != nil {
			cancel()
		}
		return err
	}
	return errors.New("service attempted to stop when not running")
}

type Service interface {
	Name() string
	IsRunning() bool
	Run() error
	Stop(timeout time.Duration) error
}
