package neo

import "net/http"

type Context interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

type Middleware func(next Handler) Handler
type SimpleMiddleware func(c Context, next Handler) error

type Chain struct {
	Middlewares []Middleware
}

type Handler interface {
	Run(Context) error
}

type HandlerFunc func(Context) error

func (f HandlerFunc) Run(c Context) error {
	return f(c)
}
