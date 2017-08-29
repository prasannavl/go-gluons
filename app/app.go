package app

import (
	"fmt"
	"net/http"
	"pvl/apicore/app/reqcontext"
	"time"

	"context"

	"github.com/go-chi/render"
	"github.com/prasannavl/go-grab/lifecycle"
	"github.com/prasannavl/go-grab/log"
	"github.com/prasannavl/mchain"
	"github.com/prasannavl/mchain/builder"
)

func createAppContext(logger *log.Logger, addr string) *AppContext {
	services := Services{Logger: logger}
	c := AppContext{Services: services, ServerAddress: addr}
	return &c
}

func newAppHandler(c *AppContext) http.Handler {
	reqLogHandler := reqcontext.CreateRequestLogHandler(c.Logger)
	b := builder.Create()
	b.AddSimple(
		reqcontext.RequestContextInitHandler,
		reqLogHandler,
		reqcontext.RequestDurationHandler,
		reqcontext.CreateReqIDHandler(false),
	)
	b.Handler(CreateActionHandler(c.ServerAddress))
	return b.BuildHttp(func(err error) {
		c.Logger.Errorf("unhandled: %s", err.Error())
	})
}

func NewApp(context *AppContext) http.Handler {
	m := http.NewServeMux()
	m.Handle("/", newAppHandler(context))
	return http.Handler(m)
}

func CreateActionHandler(host string) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		data := struct {
			Message string
			Date    time.Time
		}{
			fmt.Sprintf("Hello world from %s", host),
			time.Now(),
		}
		render.JSON(w, r, &data)
		return nil
	}
	return mchain.HandlerFunc(f)
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
