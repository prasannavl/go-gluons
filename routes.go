package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type Services struct {
	log *log.Logger
}

func NewHandlerFunc(host string, logger *log.Logger) func(http.ResponseWriter, *http.Request) {

	services := Services{log: logger}
	var _ = services

	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Message string
			Date    time.Time
		}{
			fmt.Sprintf("Hello world from %s", host),
			time.Now(),
		}
		render.JSON(w, r, &data)
	}
}
