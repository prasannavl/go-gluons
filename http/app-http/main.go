package main

import (
	"net"
	"strconv"

	flag "github.com/spf13/pflag"

	"fmt"

	"net/http"

	"github.com/prasannavl/go-gluons/http/pprof"

	"github.com/prasannavl/go-gluons/appx"
	"github.com/prasannavl/go-gluons/http/app-http/app"
	"github.com/prasannavl/go-gluons/http/utils"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/go-gluons/logconfig"
)

type EnvFlags struct {
	Addr           string
	LogFile        string
	LogDisabled    bool
	Verbosity      int
	DisplayVersion bool
	PprofAddr      string
	LogHumanize    bool
	Insecure       bool
	RedirectorAddr string
	UseSelfSigned  bool
	Hosts          []string
}

func initFlags(env *EnvFlags) {
	flag.BoolVar(&env.DisplayVersion, "version", false, "display the version and exit")
	flag.CountVarP(&env.Verbosity, "verbose", "v", "verbosity level")
	flag.StringVarP(&env.Addr, "address", "a", "localhost:8000", "the 'host:port' for the service to listen on")
	flag.StringVar(&env.PprofAddr, "pprof-address", "localhost:9090", "the 'host:port' for pprof")
	flag.StringVar(&env.LogFile, "log", "", "the log file destination")
	flag.BoolVar(&env.LogDisabled, "no-log", false, "disable the logger")
	flag.BoolVarP(&env.LogHumanize, "log-humanize", "h", false, "humanize log messages")
	flag.BoolVar(&env.Insecure, "insecure", false, "disable tls")
	flag.BoolVar(&env.UseSelfSigned, "self-signed", false, "use randomly generated self signed certificate for tls")
	flag.StringVar(&env.RedirectorAddr, "redirector", "", "a redirector address as 'host:port' to enable")
	flag.StringArrayVar(&env.Hosts, "hosts", nil, "'host:port' items to enable hosts filter")

	flag.Usage = func() {
		printPackageHeader(false)
		fmt.Printf("Usage: [opts]\n\nOptions:\r\n")
		flag.PrintDefaults()
		fmt.Println()
	}
}

func initLogging(env *EnvFlags) logconfig.LogInitResult {
	logInitResult := logconfig.LogInitResult{}
	if !env.LogDisabled {
		logOpts := logconfig.DefaultOptions()
		if !env.LogHumanize {
			logOpts.Humanize = logconfig.Humanize.False
		}
		if env.LogFile != "" {
			logOpts.LogFile = env.LogFile
		}
		logOpts.VerbosityLevel = env.Verbosity
		logconfig.Init(&logOpts, &logInitResult)
	}
	return logInitResult
}

func printPackageHeader(versionOnly bool) {
	if versionOnly {
		fmt.Printf("%s", app.Version)
	} else {
		fmt.Printf("%s\t%s\r\n", app.Package, app.Version)
	}
}

func tryRunRedirector(redirectAddr string, addr string) {
	if redirectAddr == "" {
		return
	}
	hostAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		l, err := net.Listen("tcp", redirectAddr)
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
			utils.Redirect(w, r, finalAddr, http.StatusPermanentRedirect)
		}
		http.Serve(l, http.HandlerFunc(f))
	}()
}

func main() {
	appx.InitTerm()

	env := EnvFlags{}
	initFlags(&env)

	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			log.Errorv(err)
		}
	}()

	if env.DisplayVersion {
		printPackageHeader(true)
		return
	}

	logInitResult := initLogging(&env)
	log.Infof("listen-address: %s", env.Addr)
	if env.PprofAddr != "" {
		go pprof.Run(env.PprofAddr, "/pprof")
	}
	tryRunRedirector(env.RedirectorAddr, env.Addr)
	app.Run(logInitResult.Logger, env.Addr, env.Hosts, env.Insecure, env.UseSelfSigned)
}
