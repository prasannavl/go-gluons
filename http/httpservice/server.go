package httpservice

import (
	"crypto/tls"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prasannavl/go-gluons/cert"
	"github.com/prasannavl/go-gluons/log"
	"golang.org/x/crypto/acme/autocert"
)

type HandlerServiceOpts struct {
	Logger        *log.Logger
	Addr          string
	Handler       http.Handler
	WebRoot       string
	ServiceName   string
	Insecure      bool
	Hosts         []string
	CacheDir      string
	UseSelfSigned bool
}

type TLSOpts struct {
	Addr          string
	Hosts         []string
	CacheDir      string
	UseSelfSigned bool
}

func NewHandlerService(opts *HandlerServiceOpts) (Service, error) {
	logger := opts.Logger
	addr := opts.Addr

	var stdErrLog *stdlog.Logger
	if logger != nil {
		stdErrLog = stdlog.New(log.NewLogWriter(logger, log.ErrorLevel, ""), "", 0)
	}
	server := &http.Server{
		Handler:        opts.Handler,
		IdleTimeout:    20 * time.Second,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		ErrorLog:       stdErrLog,
		MaxHeaderBytes: 1 << 20, // 1mb
	}

	var listener net.Listener

	var err error
	if opts.Insecure {
		listener, err = net.Listen("tcp", addr)
	} else {
		tlsConf := tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		}
		tlsOpts := TLSOpts{
			Addr:          addr,
			CacheDir:      opts.CacheDir,
			UseSelfSigned: opts.UseSelfSigned,
			Hosts:         opts.Hosts,
		}
		listener, err = CreateTLSListener(&tlsOpts, &tlsConf)
	}
	if err != nil {
		return nil, err
	}

	serviceName := opts.ServiceName
	if serviceName == "" {
		serviceName = addr
	}
	return New(serviceName, server, listener), nil
}

func CreateTLSListener(opts *TLSOpts, tlsConf *tls.Config) (net.Listener, error) {
	tlsHosts := opts.Hosts
	if opts.UseSelfSigned || len(tlsHosts) == 0 {
		tcert, err := cert.CreateSelfSignedRandomX509("Local", nil)
		if err != nil {
			return nil, err
		}
		tlsConf.Certificates = []tls.Certificate{tcert}
	} else {
		certMgr := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(tlsHosts...),
			Cache:      autocert.DirCache(opts.CacheDir),
		}
		tlsConf.GetCertificate = certMgr.GetCertificate
	}
	return tls.Listen("tcp", opts.Addr, tlsConf)
}

func CreateHandlerServiceOpts(addr string, handler http.Handler) HandlerServiceOpts {
	wd, _ := os.Getwd()
	return HandlerServiceOpts{
		Addr:     addr,
		Handler:  handler,
		Insecure: true,
		WebRoot:  wd,
	}
}

func CreateSelfSignedTLSOpts(addr string) TLSOpts {
	return TLSOpts{
		Addr:          addr,
		UseSelfSigned: true,
	}
}

func CreateTLSOpts(addr string, certCacheDir string, hosts ...string) TLSOpts {
	return TLSOpts{
		Addr:     addr,
		CacheDir: certCacheDir,
		Hosts:    hosts,
	}
}
