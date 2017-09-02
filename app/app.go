package app

import (
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/prasannavl/go-grab/lifecycle"
	"github.com/prasannavl/go-grab/log"
	"github.com/prasannavl/go-httpapi-base/app/reqcontext"
	"github.com/prasannavl/go-httpapi-base/app/responder"
	"github.com/prasannavl/mchain/builder"
)

func createAppContext(logger *log.Logger, addr string) *AppContext {
	services := Services{Logger: logger}
	c := AppContext{Services: services, ServerAddress: addr}
	return &c
}

func newAppHandler(c *AppContext) http.Handler {
	b := builder.CreateHttp()

	b.Add(
		reqcontext.CreateInitHandler(c.Logger),
		reqcontext.LogHandler,
		reqcontext.CreateRequestIDHandler(false),
	)

	b.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendReply(w, r)
	}))

	return b.Build()
}

func NewApp(context *AppContext) http.Handler {
	m := http.NewServeMux()
	m.Handle("/", newAppHandler(context))
	return http.Handler(m)
}

func sendReply(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Message string
		Date    time.Time
	}{
		fmt.Sprint("Hello world!"),
		time.Now(),
	}
	responder.Send(w, r, &data)
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
