package diag

import (
	"net/http"
	"strings"

	"github.com/prasannavl/mchain/hconv"

	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/utils"
	"github.com/prasannavl/go-gluons/log"
)

func LogLevelSwitcher(path string, paramName string) func(*http.ServeMux) {
	return func(mux *http.ServeMux) {
		mux.HandleFunc(path, hconv.FuncToHttp(LogLevelSwitchHandlerFunc(paramName), utils.LoggedHttpCodeOrInternalServerError))
	}
}

func LogLevelSwitchHandlerFunc(paramName string) mchain.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) error {
		lvl := r.URL.Query().Get(paramName)
		if lvl != "" {
			lvl = strings.ToLower(strings.TrimSpace(lvl))
			level := log.LogLevelFromString(lvl)
			if log.IsValidLevel(level) {
				log.SetFilter(log.GetLogger(), log.LogFilterForLevel(level))
				w.WriteHeader(http.StatusOK)
				return nil
			}
		}
		return httperror.New(http.StatusBadRequest, "invalid log level", true)
	}
	return f
}
