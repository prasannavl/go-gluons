package app

import (
	"fmt"
	"net/http"
	"time"

	"context"

	stdlog "log"

	"github.com/prasannavl/go-gluons/httputils/app/reqcontext"
	"github.com/prasannavl/go-gluons/httputils/app/responder"
	"github.com/prasannavl/go-gluons/lifecycle"
	"github.com/prasannavl/go-gluons/log"
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
		reqcontext.CreateLogHandler(log.InfoLevel),
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

	stdErrLog := stdlog.New(log.NewLogWriter(logger, log.ErrorLevel, ""), "", 0)
	server := &http.Server{
		Addr:           c.ServerAddress,
		Handler:        a,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		ErrorLog:       stdErrLog,
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
