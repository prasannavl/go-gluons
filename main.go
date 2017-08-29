package main

import (
	flag "github.com/spf13/pflag"

	"fmt"

	"pvl/apicore/app"
	"pvl/apicore/platform"

	"github.com/prasannavl/go-grab/log"
	logc "github.com/prasannavl/go-grab/log-config"
)

func main() {
	var addr string
	var logFile string
	var logDisabled bool
	var verbosity int
	var displayVersion bool

	platform.Init()

	flag.Usage = func() {
		printPackageHeader(false)
		fmt.Printf("Usage: [opts]\n\nOptions:\r\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.BoolVar(&displayVersion, "version", false, "display the version and exit")
	flag.CountVarP(&verbosity, "verbose", "v", "verbosity level")
	flag.StringVarP(&addr, "address", "a", "localhost:8000", "the 'host:port' for the service to listen on")
	flag.StringVar(&logFile, "log", "", "the log file destination")
	flag.BoolVar(&logDisabled, "no-log", false, "disable the logger")

	flag.Parse()

	if displayVersion {
		printPackageHeader(true)
		return
	}

	logInitResult := logc.LogInitResult{}
	if !logDisabled {
		logOpts := logc.DefaultOptions()
		if logFile != "" {
			logOpts.LogFile = logFile
		}
		logOpts.VerbosityLevel = verbosity
		logc.Init(&logOpts, &logInitResult)
	}

	log.Infof("listen-address: %q", addr)
	app.Run(logInitResult.Logger, addr)
}

func printPackageHeader(versionOnly bool) {
	if versionOnly {
		fmt.Printf("%s", app.Version)
	} else {
		fmt.Printf("%s\t%s\r\n", app.Package, app.Version)
	}
}
