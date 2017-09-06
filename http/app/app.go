package app

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"context"

	stdlog "log"

	"github.com/prasannavl/go-gluons/appx"
	"github.com/prasannavl/go-gluons/http/chainutils"
	"github.com/prasannavl/go-gluons/http/fileserver"
	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/go-gluons/http/reqcontext"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/mchain/builder"
	"github.com/prasannavl/mchain/hconv"
	"github.com/rsms/gotalk"
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
	b := builder.Create()

	// api := gotalk.NewHandlers()

	b.Add(
		reqcontext.CreateInitMiddleware(c.Logger),
		reqcontext.CreateLogMiddleware(log.InfoLevel),
		middleware.ErrorHandlerMiddleware,
		middleware.PanicRecoveryMiddleware,
		reqcontext.CreateRequestIDHandler(false),
		chainutils.Mount("/api", hconv.FromHttp(gotalk.WebSocketHandler())),
	)

	b.Handler(fileserver.NewEx(http.Dir("./app/www"), NotFoundHandler))
	return b.BuildHttp(nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) error {
	f, err := os.Open("./app/www/error/404.html")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	return nil
}

func NewApp(context *AppContext) http.Handler {
	m := http.NewServeMux()
	m.Handle("/", newAppHandler(context))
	return http.Handler(m)
}

func Run(logger *log.Logger, addr string) {
	c := createAppContext(logger, addr)
	a := NewApp(c)

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

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Errorf("server close: %s", err.Error())
	}
	log.Info("exit")
}
