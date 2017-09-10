package app

import (
	"crypto/tls"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prasannavl/mchain/hconv"

	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/cert"
	"github.com/prasannavl/go-gluons/http/chainutils"
	"github.com/prasannavl/go-gluons/http/gosock"
	"github.com/prasannavl/go-gluons/http/hostrouter"
	"github.com/prasannavl/go-gluons/http/utils"
	"golang.org/x/crypto/acme/autocert"

	"context"

	stdlog "log"

	"github.com/prasannavl/go-gluons/appx"
	"github.com/prasannavl/go-gluons/http/fileserver"
	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/mchain/builder"
)

func createAppContext(logger *log.Logger, addr string) *AppContext {
	services := Services{
		Logger:        logger,
		TemplateCache: make(map[string]*template.Template),
	}
	c := AppContext{
		Services:      services,
		ServerAddress: addr,
	}
	return &c
}

func newAppHandler(c *AppContext, webRoot string) mchain.Handler {
	apiHandlers := apiHandlers(c)
	wss := gosock.NewWebSocketServer(apiHandlers)

	b := builder.Create()

	b.Add(
		middleware.CreateInitMiddleware(c.Logger),
		middleware.CreateLogMiddleware(log.InfoLevel),
		middleware.ErrorHandlerMiddleware,
		middleware.PanicRecoveryMiddleware,
		middleware.CreateRequestIDHandler(false),
		chainutils.OnPrefix("/api", wss),
		chainutils.OnPrefix("/assets/gotalk.js", gosock.CreateAssetHandler("/assets/gotalk.js", "/api", false)),
	)

	b.Handler(fileserver.NewEx(http.Dir(webRoot), CreateNotFoundHandler(webRoot).ServeHTTP))
	return b.Build()
}

func CreateNotFoundHandler(webRoot string) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		f, err := os.Open(webRoot + "/error/404.html")
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusNotFound)
		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
		return nil
	}
	return mchain.HandlerFunc(f)
}

func NewApp(context *AppContext, webRoot string, hosts []string) http.Handler {
	appHandler := hconv.ToHttp(newAppHandler(context, webRoot), nil)
	if len(hosts) == 0 {
		return appHandler
	}
	r := hostrouter.New()
	log.Infof("host filters: %v", hosts)
	for _, h := range hosts {
		r.HandlePattern(h, appHandler)
	}
	return r.Build(hconv.FuncToHttp(
		CreateNotFoundHandler(webRoot).ServeHTTP,
		utils.InternalServerErrorAndLog))
}

func Run(logger *log.Logger, addr string, webRoot string, hosts []string, insecure bool, useSelfSigned bool) {
	c := createAppContext(logger, addr)
	a := NewApp(c, webRoot, hosts)

	stdErrLog := stdlog.New(log.NewLogWriter(logger, log.ErrorLevel, ""), "", 0)
	server := &http.Server{
		Addr:           c.ServerAddress,
		Handler:        a,
		IdleTimeout:    20 * time.Second,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		ErrorLog:       stdErrLog,
		MaxHeaderBytes: 1 << 20, // 1mb
	}

	appx.CreateShutdownHandler(func() {
		server.Shutdown(context.Background())
	}, appx.ShutdownSignals...)

	var listener net.Listener

	var err error
	if insecure {
		listener, err = net.Listen("tcp", c.ServerAddress)
	} else {
		tlsConf := tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		}
		listener, err = createTLSListener(addr, &tlsConf, useSelfSigned, hosts)
	}
	if err != nil {
		log.Errorf("server-listener: %v", err)
		return
	}
	err = server.Serve(listener)
	if err != http.ErrServerClosed {
		log.Errorf("server-close: %v", err)
	}
	log.Info("exit")
}

func createTLSListener(addr string, tlsConf *tls.Config, useSelfSigned bool, tlsHosts []string) (net.Listener, error) {
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
