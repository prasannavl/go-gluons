package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type httpHandler func(http.ResponseWriter, *http.Request)

func main() {
	var addr string
	var logFile string
	var logOverwrite bool
	var logDisable bool
	var debugMode bool

	flag.Usage = func() {
		fmt.Printf("\nUsage: [opts]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.StringVar(&addr, "address", "localhost:8000", "the `address:port` for the service to listen on")
	flag.StringVar(&logFile, "log", "", "the log file destination")
	flag.BoolVar(&logOverwrite, "log-overwrite", false, "overwrite the log file if set, or appends")
	flag.BoolVar(&logDisable, "log-disable", false, "disable the logger")
	flag.Parse()

	log.Printf("Listen address: %s", addr)

	if logDisable {
		log.SetOutput(ioutil.Discard)
	} else {
		logWriter := CreateLogStream(&logFile, logOverwrite)
		defer logWriter.Close()
		log.SetOutput(logWriter)
	}

	runServer(addr, handler(addr))
}

func NewHandler(host string) httpHandler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		data := struct {
			Message string
			Date    time.Time
		}{
			fmt.Sprintf("Hello world from %s", host),
			time.Now(),
		}

		render.JSON(w, r, &data)
	})
	return r.ServeHTTP
}

func handler(addr string) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("labs.prasannavl.com/", NewHandler("PVL Labs"))
	m.HandleFunc(addr+"/", NewHandler(addr))
	m.HandleFunc("nf.statwick.com/", NewHandler("NextFirst API"))
	return m
}

func runServer(addr string, handler http.Handler) {
	server := &http.Server{Addr: addr, Handler: handler}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for sig := range quit {
			log.Printf("Signal: %q", sig)
			log.Printf("Shutting down..")
			server.Shutdown(context.Background())
		}
	}()

	log.Println("Listening..")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Fatal: %v", err)
	}
	log.Println("Exit: Success")
}
