package main

import (
	"net/http"
	"path/filepath"
	"runtime/debug"

	flag "github.com/spf13/pflag"

	"fmt"

	"github.com/prasannavl/go-gluons/http/diag"
	"github.com/prasannavl/go-gluons/http/httpservice"
	"github.com/prasannavl/go-gluons/http/redirector"

	"github.com/prasannavl/go-gluons/appx"
	"github.com/prasannavl/go-gluons/http/app-http/app"
	"github.com/prasannavl/go-gluons/log"
	"github.com/prasannavl/go-gluons/logconfig"
)

type EnvFlags struct {
	Addr           string
	LogFile        string
	LogDisabled    bool
	Verbosity      int
	DisplayVersion bool
	DiagAddr       string
	LogHumanize    bool
	Insecure       bool
	RedirectorAddr string
	UseSelfSigned  bool
	Hosts          []string
	WebRoot        string
	CertCacheDir   string
}

func initFlags(env *EnvFlags) {
	flag.BoolVar(&env.DisplayVersion, "version", false, "display the version and exit")
	flag.CountVarP(&env.Verbosity, "verbose", "v", "verbosity level")
	flag.StringVarP(&env.Addr, "address", "a", "localhost:8000", "the 'host:port' for the service to listen on")
	flag.StringVar(&env.DiagAddr, "dapi-address", "localhost:9090", "the 'host:port' for diagnostics api")
	flag.StringVar(&env.LogFile, "log", "", "the log file destination")
	flag.BoolVar(&env.LogDisabled, "no-log", false, "disable the logger")
	flag.BoolVarP(&env.LogHumanize, "log-humanize", "h", false, "humanize log messages")
	flag.BoolVar(&env.Insecure, "insecure", false, "disable tls")
	flag.BoolVar(&env.UseSelfSigned, "self-signed", false, "use randomly generated self signed certificate for tls")
	flag.StringVar(&env.RedirectorAddr, "redirector", "", "a redirector address as 'host:port' to enable")
	flag.StringArrayVar(&env.Hosts, "hosts", nil, "'host:port' items to enable hosts filter")
	flag.StringVar(&env.WebRoot, "root", "", "web root path")
	flag.StringVar(&env.CertCacheDir, "cert-dir", "", "the auto-tls certificate cache dir")

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

func main() {
	appx.InitTerm()

	env := EnvFlags{}
	initFlags(&env)

	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			log.Errorv(err)
			debug.PrintStack()
		}
	}()

	if env.DisplayVersion {
		printPackageHeader(true)
		return
	}

	logInitResult := initLogging(&env)
	log.Infof("listen-address: %s", env.Addr)

	if env.DiagAddr != "" {
		s1 := diag.CreateWithConfigure(env.DiagAddr,
			diag.SetupIndexNotFound,
			diag.LogLevelSwitcher(nil))
		go s1.Run()
	}

	if env.RedirectorAddr != "" {
		s2 := redirector.Create(env.RedirectorAddr, env.Addr)
		go s2.Run()
	}

	opts := httpservice.HandlerServiceOpts{
		Addr:          env.Addr,
		Logger:        logInitResult.Logger,
		WebRoot:       filepath.Clean(env.WebRoot),
		Hosts:         env.Hosts,
		CacheDir:      env.CertCacheDir,
		UseSelfSigned: env.UseSelfSigned,
	}

	if env.Insecure {
		opts.Insecure = true
	}

	service, err := app.CreateService(&opts)

	if err != nil {
		log.Errorf("failed to create service: %v", err)
		return
	}

	appx.CreateShutdownHandler(func() {
		_ = service.Stop(0)
	}, appx.ShutdownSignals...)

	err = service.Run()
	if err != http.ErrServerClosed {
		log.Errorf("service: %v", err)
	}

	log.Info("exit")
}
