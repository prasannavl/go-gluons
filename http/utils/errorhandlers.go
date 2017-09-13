package utils

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/log"
)

var HttpCodeOrInternalServerError = CreateHttpErrorHandler(http.StatusInternalServerError)
var LoggedHttpCodeOrInternalServerError = CreateHttpErrorHandler(http.StatusInternalServerError)

var InternalServerError = CreateStatusErrorHandler(http.StatusInternalServerError)
var LoggedInternalServerError = CreateLoggedStatusErrorHandler(http.StatusInternalServerError)

var BadRequestError = CreateStatusErrorHandler(http.StatusBadRequest)
var LoggedBadRequestError = CreateLoggedStatusErrorHandler(http.StatusBadRequest)

func CreateLoggedStatusErrorHandler(status int) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		log.Errorf("error-handler: %v", err)
		w.WriteHeader(status)
	}
}

func CreateStatusErrorHandler(status int) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
}

func CreateLoggedHttpErrorHandler(fallbackStatus int) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		if e, ok := err.(httperror.HttpError); ok {
			if httperror.IsServerErrorCode(e.Code()) {
				log.Errorf("error-handler: %v", e)
			}
			w.WriteHeader(e.Code())
		} else {
			log.Errorf("error-handler: %v", err)
			w.WriteHeader(fallbackStatus)
		}
	}
}

func CreateHttpErrorHandler(fallbackStatus int) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		if e, ok := err.(httperror.HttpError); ok {
			w.WriteHeader(e.Code())
		} else {
			w.WriteHeader(fallbackStatus)
		}
	}
}
