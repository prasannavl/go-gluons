package hconv

import (
	"net/http"
	"pvl/apicore/app/neo"
)

func FuncToHttpFactory(h neo.HandlerFunc, errorHandler func(error)) func(neo.Context) http.HandlerFunc {
	hh := func(c neo.Context) http.HandlerFunc {
		h := func(w http.ResponseWriter, r *http.Request) {
			err := h.Run(c)
			if err != nil && errorHandler != nil {
				errorHandler(err)
			}
		}
		return http.HandlerFunc(h)
	}
	return hh
}

func ToHttpFactory(h neo.Handler, errorHandler func(error)) func(neo.Context) http.Handler {
	hf := neo.HandlerFunc(h.Run)
	hh := FuncToHttpFactory(hf, errorHandler)
	hhf := func(c neo.Context) http.Handler {
		return http.Handler(hh(c))
	}
	return hhf
}

func FuncFromHttp(h http.HandlerFunc) neo.HandlerFunc {
	hh := func(c neo.Context) error {
		h.ServeHTTP(c.ResponseWriter(), c.Request())
		return nil
	}
	return neo.HandlerFunc(hh)
}

func FuncFromHttpRecoverable(h http.HandlerFunc) neo.HandlerFunc {
	hh := func(c neo.Context) error {
		var err error
		defer neo.RecoverIntoError(&err)
		h.ServeHTTP(c.ResponseWriter(), c.Request())
		return err
	}
	return neo.HandlerFunc(hh)
}

func FromHttp(h http.Handler) neo.Handler {
	hh := FuncFromHttp(http.HandlerFunc(h.ServeHTTP))
	return neo.Handler(hh)
}

func FromHttpRecoverable(h http.Handler) neo.Handler {
	hh := FuncFromHttpRecoverable(http.HandlerFunc(h.ServeHTTP))
	return neo.Handler(hh)
}
