package pprof

import (
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/prasannavl/go-gluons/log"
)

func Run(addr string, path string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("pprof-listener: %v", err)
	}
	log.Infof("pprof endpoint: %s", l.Addr())
	mux := http.NewServeMux()
	SetupHandlers(mux, path)
	http.Serve(l, mux)
}

func SetupHandlers(mux *http.ServeMux, pathPrefix string) {
	mux.Handle(pathPrefix, http.HandlerFunc(pprof.Index))
	mux.Handle(pathPrefix+"/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle(pathPrefix+"/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle(pathPrefix+"/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle(pathPrefix+"/trace", http.HandlerFunc(pprof.Trace))
}
