package pprof

import (
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prasannavl/go-gluons/http/httpservice"
	"github.com/prasannavl/go-gluons/log"
)

func SetupHandlers(mux *http.ServeMux, pathPrefix string) {
	mux.Handle(pathPrefix, http.HandlerFunc(pprof.Index))
	mux.Handle(pathPrefix+"/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle(pathPrefix+"/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle(pathPrefix+"/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle(pathPrefix+"/trace", http.HandlerFunc(pprof.Trace))
}

func Create(addr string, path string) httpservice.Service {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("pprof-listener: %v", err)
	}
	log.Infof("pprof endpoint: %s", l.Addr())
	mux := http.NewServeMux()
	SetupHandlers(mux, path)
	server := &http.Server{
		Handler:        mux,
		IdleTimeout:    20 * time.Second,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   60 * time.Second,
		ErrorLog:       nil,
		MaxHeaderBytes: 1 << 12, // 4kb
	}
	return httpservice.New("pprof", server, l)
}
