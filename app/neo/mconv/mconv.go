package mconv

import (
	"net/http"
	"pvl/apicore/app/neo"
	"pvl/apicore/app/neo/hconv"
)

func FromSimple(fn func(c neo.Context, next neo.Handler) error) neo.Middleware {
	m := func(next neo.Handler) neo.Handler {
		f := func(ctx neo.Context) error {
			return fn(ctx, next)
		}
		return neo.HandlerFunc(f)
	}
	return m
}

func ToSimple(middleware neo.Middleware) (fn func(c neo.Context, next neo.Handler) error) {
	h := func(c neo.Context, next neo.Handler) error {
		return middleware(next).Run(c)
	}
	return h
}

func FromHttp(h func(http.Handler) http.Handler, innerErrorHandler func(error)) neo.Middleware {
	hh := func(hx neo.Handler) neo.Handler {
		httpHandlerFactory := hconv.ToHttpFactory(hx, innerErrorHandler)
		nh := func(context neo.Context) error {
			nextHandler := httpHandlerFactory(context)
			h(nextHandler).ServeHTTP(context.ResponseWriter(), context.Request())
			return nil
		}
		return neo.HandlerFunc(nh)
	}
	return neo.Middleware(hh)
}

func FromHttpRecoverable(h func(http.Handler) http.Handler, innerErrorHandler func(error)) neo.Middleware {
	hh := func(hx neo.Handler) neo.Handler {
		httpHandlerFactory := hconv.ToHttpFactory(hx, innerErrorHandler)
		nh := func(context neo.Context) error {
			var err error
			defer neo.RecoverIntoError(&err)
			nextHandler := httpHandlerFactory(context)
			h(nextHandler).ServeHTTP(context.ResponseWriter(), context.Request())
			return err
		}
		return neo.HandlerFunc(nh)
	}
	return neo.Middleware(hh)
}
