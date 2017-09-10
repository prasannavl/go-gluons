package httpservice

import (
	"crypto/tls"
	stdlog "log"
	"net"
	"net/http"
	"time"

	"github.com/prasannavl/go-gluons/cert"
	"github.com/prasannavl/go-gluons/log"
	"golang.org/x/crypto/acme/autocert"
)

func NewHandlerService(logger *log.Logger, addr string, handler http.Handler,
	webRoot string, hosts []string, insecure bool, useSelfSigned bool) (Service, error) {

	stdErrLog := stdlog.New(log.NewLogWriter(logger, log.ErrorLevel, ""), "", 0)
	server := &http.Server{
		Handler:        handler,
		IdleTimeout:    20 * time.Second,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		ErrorLog:       stdErrLog,
		MaxHeaderBytes: 1 << 20, // 1mb
	}

	var listener net.Listener

	var err error
	if insecure {
		listener, err = net.Listen("tcp", addr)
	} else {
		tlsConf := tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		}
		listener, err = CreateTLSListener(addr, &tlsConf, useSelfSigned, hosts)
	}
	if err != nil {
		return nil, err
	}

	return New(addr, server, listener), nil
}

func CreateTLSListener(addr string, tlsConf *tls.Config, useSelfSigned bool, tlsHosts []string) (net.Listener, error) {
	if useSelfSigned || len(tlsHosts) == 0 {
		tcert, err := cert.CreateSelfSignedRandomX509("Local", nil)
		if err != nil {
			return nil, err
		}
		tlsConf.Certificates = []tls.Certificate{tcert}
	} else {
		certMgr := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(tlsHosts...),
		}
		tlsConf.GetCertificate = certMgr.GetCertificate
	}
	return tls.Listen("tcp", addr, tlsConf)
}
