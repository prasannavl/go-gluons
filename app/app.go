package app

import (
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/prasannavl/go-grab/lifecycle"
	"github.com/prasannavl/go-grab/log"
	"github.com/prasannavl/go-starter-api/app/middleware"
	"github.com/prasannavl/go-starter-api/app/reqcontext"
	"github.com/prasannavl/go-starter-api/app/responder"
	"github.com/prasannavl/mchain"
	"github.com/prasannavl/mchain/builder"
)

func createAppContext(logger *log.Logger, addr string) *AppContext {
	services := Services{Logger: logger}
	c := AppContext{Services: services, ServerAddress: addr}
	return &c
}

func newAppHandler(c *AppContext) http.Handler {
	b := builder.Create()

	b.Add(
		reqcontext.CreateInitHandler(c.Logger),
		reqcontext.ErrorHandler,
		reqcontext.LogHandler,
		reqcontext.DurationHandler,
		middleware.RecoverPanicHandler,
		reqcontext.CreateRequestIDHandler(false),
	)

	b.Handler(CreateActionHandler(c.ServerAddress))
	return b.BuildHttp(func(err error) {
		c.Logger.Errorf("unhandled: %#v", err)
	})
}

func NewApp(context *AppContext) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		sendReply(context.ServerAddress, w, r)
	})
	m.Handle("/", newAppHandler(context))
	return http.Handler(m)
}

func CreateActionHandler(host string) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		sendReply(host, w, r)
		return nil
	}
	return mchain.HandlerFunc(f)
}

func sendReply(host string, w http.ResponseWriter, r *http.Request) {
	data := struct {
		Message string
		Date    time.Time
	}{
		fmt.Sprintf("Hello world from %s", host),
		time.Now(),
	}
	responder.Send(&data, w, r)
}

func Run(logger *log.Logger, addr string) {
	c := createAppContext(logger, addr)
	a := NewApp(c)
	server := &http.Server{
		Addr:           c.ServerAddress,
		Handler:        a,
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
