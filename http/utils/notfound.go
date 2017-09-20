package utils

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"
)

// NotFound replies to the request with an HTTP 404 not found error.
func HttpNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// NotFoundHandler returns a simple request handler
// that replies to each request with a ``404 page not found'' reply.
func HttpNotFoundHandler() http.Handler { return http.HandlerFunc(HttpNotFound) }

func NotFoundToError(w http.ResponseWriter, r *http.Request) error {
	return httperror.New(http.StatusNotFound, "", true)
}

func NotFoundToErrorHandler() mchain.Handler {
	return mchain.HandlerFunc(NotFoundToError)
}

func HttpNotFoundText(w http.ResponseWriter, r *http.Request) {
	HttpNotFound(w, r)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func HttpNotFoundTextHandler() http.Handler { return http.HandlerFunc(HttpNotFoundText) }
