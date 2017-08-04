package main

import (
	"fmt"
	"net/http"
	"time"

	"pvl/api-core/middleware"

	"github.com/go-chi/render"
	"github.com/prasannavl/mchain"
)

func NewApp(addr string) http.Handler {
	m := http.NewServeMux()
	m.Handle("labs.prasannavl.com/", newAppHandler("PVL Labs"))
	m.Handle(addr+"/", newAppHandler(addr))
	m.Handle("nf.statwick.com/", newAppHandler("NextFirst API"))
	return http.Handler(m)
}

func newAppHandler(host string) http.Handler {
	// c := appcontext.AppContext{Services: appcontext.Services{}}
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
