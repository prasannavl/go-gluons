package responder

import (
	"fmt"
	"net/http"

	"github.com/prasannavl/goerror/httperror"

	"github.com/go-chi/render"
)

// TODO: Proper content negotiation
// TODO: Use Content-Encoding

func Send(w http.ResponseWriter, r *http.Request, value interface{}) {
	if value == nil {
		return
	}
	render.JSON(w, r, value)
}

func SendError(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := err.(httperror.HttpError); ok {
		sendHttpError(w, r, e)
		return
	}
	SendWithStatus(w, r, http.StatusInternalServerError, err.Error())
}

func SendWithStatus(w http.ResponseWriter, r *http.Request, status int, value interface{}) {
	SetStatus(w, status)
	Send(w, r, value)
}

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

func sendHttpError(w http.ResponseWriter, r *http.Request, err httperror.HttpError) {
	msg := err.Error()
	SendWithStatus(w, r, err.Code(), msg)
}
