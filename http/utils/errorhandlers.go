package utils

import (
	"net/http"

	"github.com/prasannavl/go-gluons/log"
)

func InternalServerError(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func InternalServerErrorAndLog(err error, w http.ResponseWriter, r *http.Request) {
	log.Errorf("error-handler: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}
