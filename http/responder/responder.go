package responder

import (
	"fmt"
	"net/http"

	"github.com/prasannavl/go-errors/httperror"
	"github.com/unrolled/render"
)

// TODO: Proper content negotiation
// TODO: Use Content-Encoding

// Just use a singleton for now.
var Renderer = render.New(render.Options{})

func SetStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func SendErrorText(w http.ResponseWriter, errOrStringer interface{}) {
	var code int
	var message string
	switch e := errOrStringer.(type) {
	case error:
		message = e.Error()
		if e, ok := e.(httperror.HttpError); ok {
			code = e.Code()
		}
	case string:
		message = e
	case fmt.Stringer:
		message = e.String()
	}
	c := httperror.ErrorCode(code)
	if message == "" {
		SetStatus(w, c)
	} else {
		http.Error(w, message, c)
	}
}

func Send(w http.ResponseWriter, r *http.Request, value interface{}) error {
	return SendWithStatus(w, r, http.StatusOK, value)
}

func SendWithStatus(w http.ResponseWriter, r *http.Request, status int, value interface{}) error {
	return Renderer.JSON(w, status, value)
}

func SendError(w http.ResponseWriter, r *http.Request, err error) error {
	if e, ok := err.(httperror.HttpError); ok {
		return sendHttpError(w, r, e)
	}
	return SendWithStatus(w, r, http.StatusInternalServerError, err.Error())
}

func sendHttpError(w http.ResponseWriter, r *http.Request, err httperror.HttpError) error {
	msg := err.Error()
	return SendWithStatus(w, r, err.Code(), msg)
}
