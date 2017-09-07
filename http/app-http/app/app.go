package app

import (
	"crypto/tls"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/prasannavl/go-gluons/http/chainutils"
	"github.com/prasannavl/go-gluons/http/gosock"

	"context"

	stdlog "log"

	"github.com/prasannavl/go-gluons/appx"
	"github.com/prasannavl/go-gluons/cert"
	"github.com/prasannavl/go-gluons/http/fileserver"
	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/go-gluons/http/reqcontext"
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

func newAppHandler(c *AppContext) http.Handler {
	apiHandlers := apiHandlers(c)
	wss := gosock.NewWebSocketServer(apiHandlers)

	b := builder.Create()

	b.Add(
		reqcontext.CreateInitMiddleware(c.Logger),
		reqcontext.CreateLogMiddleware(log.InfoLevel),
		middleware.ErrorHandlerMiddleware,
		middleware.PanicRecoveryMiddleware,
		reqcontext.CreateRequestIDHandler(false),
		chainutils.OnPrefix("/api", wss),
		chainutils.OnPrefix("/assets/gotalk.js", gosock.CreateAssetHandler("/assets/gotalk.js", "/api", false)),
	)

	b.Handler(fileserver.NewEx(http.Dir("./www"), NotFoundHandler))
	return b.BuildHttp(nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) error {
	f, err := os.Open("./www/error/404.html")
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

func NewApp(context *AppContext) http.Handler {
	return newAppHandler(context)
}

func Run(logger *log.Logger, addr string) {
	c := createAppContext(logger, addr)
	a := newAppHandler(c)

	stdErrLog := stdlog.New(log.NewLogWriter(logger, log.ErrorLevel, ""), "", 0)
	server := &http.Server{
		Addr:           c.ServerAddress,
		Handler:        a,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		ErrorLog:       stdErrLog,
		MaxHeaderBytes: 1 << 20}

	appx.CreateShutdownHandler(func() {
		server.Shutdown(context.Background())
	}, appx.ShutdownSignals...)

	tcert, _ := cert.CreateSelfSignedRandomX509("PVL Labs", nil)

	tlsConf := tls.Config{
		Certificates: []tls.Certificate{tcert},
	}

	server.TLSConfig = &tlsConf

	lsr, err := tls.Listen("tcp", c.ServerAddress, &tlsConf)
	if err != nil {
		log.Errorf("server listen: %v", err)
	}
	err = server.Serve(lsr)
	if err != http.ErrServerClosed {
		log.Errorf("server close: %v", err)
	}
	log.Info("exit")
}
