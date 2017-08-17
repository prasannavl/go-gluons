package main

import (
	"fmt"
	"net/http"
	"pvl/apicore/appcontext"
	"pvl/apicore/middleware"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/render"
	"github.com/prasannavl/mchain"
)

func NewApp(context *appcontext.AppContext) http.Handler {
	m := http.NewServeMux()
	m.Handle("labs.prasannavl.com/", newAppHandler(context, "PVL Labs"))
	m.Handle(context.ServerAddress+"/", newAppHandler(context, context.ServerAddress))
	m.Handle("nf.statwick.com/", newAppHandler(context, "NextFirst API"))
	return http.Handler(m)
}

func newAppHandler(c *appcontext.AppContext, host string) http.Handler {
	return mchain.NewBuilder(
		middleware.RequestContextInitHandler,
		middleware.RequestLogHandler,
		middleware.RequestDurationHandler,
		middleware.RequestIDMustInitHandler,
	).Handler(CreateActionHandler(host)).BuildHttp(nil)
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

func createAppContext(logger *zap.Logger, addr string) *appcontext.AppContext {
	services := appcontext.Services{Logger: logger}
	c := appcontext.AppContext{Services: services, ServerAddress: addr}
	return &c
}
