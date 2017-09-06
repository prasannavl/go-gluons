package app

import (
	"html/template"

	"github.com/prasannavl/go-gluons/log"
)

type AppContext struct {
	Services
	ServerAddress string
}

// Services is the global services context
type Services struct {
	Logger        *log.Logger
	TemplateCache map[string]*template.Template
}
