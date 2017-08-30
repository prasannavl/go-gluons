package responder

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"

	"github.com/go-chi/render"
)

// TODO: Proper content negotiation
// TODO: Use Content-Encoding

func Send(value interface{}, w http.ResponseWriter, r *http.Request) {
	if value == nil {
		return
	}
	render.JSON(w, r, value)
}

func SendContent(contentType string, value interface{}, w http.ResponseWriter, r *http.Request) {
	Send(value, w, r)
}

func SendHttpError(err httperror.Error, w http.ResponseWriter, r *http.Request) {
	var code int
	var message string
	if err != nil {
		code = httperror.ErrorCode(err.Code())
		message = err.Error()
	} else {
		code = http.StatusInternalServerError
	}
	SendWithStatus(code, message, w, r)
}

func SendError(err error, w http.ResponseWriter, r *http.Request) {
	if e, ok := err.(httperror.Error); ok {
		SendHttpError(e, w, r)
		return
	}
	if err != nil {
		Send(err.Error(), w, r)
	}
}

func SendWithStatus(status int, value interface{}, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	Send(value, w, r)
}

func SendStatus(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
}
