package diag

import (
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prasannavl/go-gluons/http/httpservice"
	"github.com/prasannavl/go-gluons/log"
)

func SetupPprof(mux *http.ServeMux) {
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
}

func Create(addr string) httpservice.Service {
	return CreateWithConfigure(addr, nil)
}

func CreateWithConfigure(addr string, configFn ...func(*http.ServeMux)) httpservice.Service {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Infof("diag endpoint: %s", l.Addr())
	mux := http.NewServeMux()
	SetupPprof(mux)
	for _, x := range configFn {
		x(mux)
	}
	server := &http.Server{
		Handler:        mux,
		IdleTimeout:    20 * time.Second,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   60 * time.Second,
		ErrorLog:       nil,
		MaxHeaderBytes: 1 << 12, // 4kb
	}
	return httpservice.New("diag", server, l)
}

func SetupIndexNotFound(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
