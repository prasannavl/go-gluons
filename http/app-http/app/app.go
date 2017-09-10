package app

import (
	"html/template"
	"io"
	"net/http"
	"os"

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
	return r.Build(hconv.FuncToHttp(
		CreateNotFoundHandler(webRoot).ServeHTTP,
		utils.InternalServerErrorAndLog))
}

func CreateService(logger *log.Logger, addr string, webRoot string, hosts []string,
	insecure bool, useSelfSigned bool) (httpservice.Service, error) {
	app := NewApp(logger, addr, webRoot, hosts)
	return httpservice.NewHandlerService(logger, addr, app,
		webRoot, hosts, insecure, useSelfSigned)
}
