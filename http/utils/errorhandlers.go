package utils

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/log"
)

var HttpCodeOrInternalServerError = CreateHttpErrorHandler(http.StatusInternalServerError, false)
var HttpCodeOrLoggedInternalServerError = CreateHttpErrorHandler(http.StatusInternalServerError, true)

var InternalServerError = CreateStatusErrorHandler(http.StatusInternalServerError, false)
var LoggedInternalServerError = CreateStatusErrorHandler(http.StatusInternalServerError, true)

var BadRequestError = CreateStatusErrorHandler(http.StatusBadRequest, false)
var LoggedBadRequestError = CreateStatusErrorHandler(http.StatusBadRequest, true)

func CreateStatusErrorHandler(status int, logErrors bool) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		if logErrors {
			log.Errorf("error-handler: %v", err)
		}
		w.WriteHeader(status)
	}
}

func CreateHttpErrorHandler(fallbackStatus int, logNonHttpErrors bool) mchain.ErrorHandler {
	return func(err error, w http.ResponseWriter, r *http.Request) {
		if e, ok := err.(httperror.HttpError); ok {
			if httperror.IsServerErrorCode(e.Code()) {
				log.Errorf("error-handler: %v", e)
			}
			w.WriteHeader(e.Code())
		} else {
			if logNonHttpErrors {
				log.Errorf("error-handler: %v", err)
			}
			w.WriteHeader(fallbackStatus)
		}
	}
}
