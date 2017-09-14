package app

import (
	"html/template"
	"net/http"

	"github.com/prasannavl/go-gluons/http/httpservice"

	"github.com/prasannavl/mchain/hconv"

	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/chainutils"
	"github.com/prasannavl/go-gluons/http/gosock"
	"github.com/prasannavl/go-gluons/http/hostrouter"
	"github.com/prasannavl/go-gluons/http/utils"

	"github.com/prasannavl/go-gluons/http/fileserver"
	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/mchain/builder"
)

// TODO: Router

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

	notFoundFilePath := webRoot + "/error/404.html"

	b.Handler(fileserver.NewEx(http.Dir(webRoot),
		utils.CreateFileHandler(notFoundFilePath, http.StatusNotFound).ServeHTTP))
	return b.Build()
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

	return r.Build(hconv.FuncToHttp(
		utils.CreateFileHandler(notFoundFilePath, http.StatusNotFound).ServeHTTP,
		utils.HttpCodeOrLoggedInternalServerError))
}

func CreateService(opts *httpservice.HandlerServiceOpts) (httpservice.Service, error) {
	app := NewApp(opts.Logger, opts.Addr, opts.WebRoot, opts.Hosts)
	opts.Handler = app
	return httpservice.NewHandlerService(opts)
}
