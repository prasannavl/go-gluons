package redirector

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prasannavl/go-gluons/http/httpservice"
	"github.com/prasannavl/go-gluons/http/utils"
	"github.com/prasannavl/go-gluons/log"
)

func Create(listenAddr string, targetAddr string) httpservice.Service {
	hostAddr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		panic(err)
	}
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Errorf("redirector-listener: %v", err)
	}
	log.Infof("redirector endpoint: %s", l.Addr())
	port := strconv.Itoa(hostAddr.Port)
	shouldIncludePort := true
	if hostAddr.Port == 443 {
		shouldIncludePort = false
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		finalAddr := "https://" + r.Host
		if shouldIncludePort {
			finalAddr += ":" + port
		}
		finalAddr += r.RequestURI
		w.Header().Set("connection", "close")
		utils.Redirect(w, r, finalAddr, http.StatusPermanentRedirect)
	}
	server := &http.Server{
		Handler:        http.HandlerFunc(f),
		IdleTimeout:    1 * time.Second,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   2 * time.Second,
		ErrorLog:       nil,
		MaxHeaderBytes: 1 << 12, // 4kb
	}
	return httpservice.New("redirector", server, l)
}
