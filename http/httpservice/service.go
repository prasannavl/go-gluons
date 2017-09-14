package httpservice

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type HttpService struct {
	name          string
	server        *http.Server
	l             net.Listener
	isRunningFlag int32
}

func New(name string, server *http.Server, listener net.Listener) Service {
	return &HttpService{name: name, server: server, l: listener}
}

func (r *HttpService) IsRunning() bool {
	return atomic.LoadInt32(&r.isRunningFlag) == 1
}

func (r *HttpService) Run() error {
	if atomic.CompareAndSwapInt32(&r.isRunningFlag, 0, 1) {
		err := r.server.Serve(r.l)
		atomic.StoreInt32(&r.isRunningFlag, 0)
		return err
	}
	return errors.New("service attempted to start while running/in-transition")
}

func (r *HttpService) Name() string {
	return r.name
}

func (r *HttpService) Stop(timeout time.Duration) error {
	if atomic.LoadInt32(&r.isRunningFlag) == 1 {
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
