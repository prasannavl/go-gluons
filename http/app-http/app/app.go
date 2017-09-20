package app

import (
	"html/template"
	"net/http"

	"github.com/prasannavl/go-gluons/http/fileserver"

	"github.com/prasannavl/go-gluons/http/httpservice"
	"github.com/prasannavl/mroute"
	"github.com/prasannavl/mroute/pat"

	"github.com/prasannavl/mchain/hconv"

	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/hostrouter"
	"github.com/prasannavl/go-gluons/http/utils"

	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/go-gluons/log"
)

// TODO: Router

func newAppHandler(c *AppContext, webRoot string) mchain.Handler {

	router := mroute.NewMux()
	dir := http.Dir(webRoot)
	router.Handle(pat.New("/*"), fileserver.NewEx(dir, nil))

	router.Use(
		middleware.CreateInitMiddleware(c.Logger),
		middleware.CreateLogMiddleware(log.InfoLevel),
		middleware.ErrorHandlerMiddleware,
		middleware.PanicRecoveryMiddleware,
		middleware.CreateRequestIDHandler(false),
	)

	return router
}

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

func NewApp(logger *log.Logger, addr string, webRoot string, hosts []string) http.Handler {
	context := createAppContext(logger, addr)
	appHandler := hconv.ToHttp(newAppHandler(context, webRoot), nil)
	if len(hosts) == 0 {
		return appHandler
	}
	r := hostrouter.New()
	log.Infof("host filters: %v", hosts)
	for _, h := range hosts {
		r.HandlePattern(h, appHandler)
	}

	notFoundFilePath := webRoot + "/error/404.html"

	return r.Build(hconv.ToHttp(
		utils.CreateFileHandler(notFoundFilePath, http.StatusNotFound),
		utils.LoggedHttpCodeOrInternalServerError))
}

func CreateService(opts *httpservice.HandlerServiceOpts) (httpservice.Service, error) {
	app := NewApp(opts.Logger, opts.Addr, opts.WebRoot, opts.Hosts)
	opts.Handler = app
	return httpservice.NewHandlerService(opts)
}
