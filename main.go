package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/net/http2"
)

var (
	config  *Config
	ch      chan<- string
	domains = make(map[string]http.Handler)

	srv, srvTLS http.Server

	confile string
	logfile string
	useTLS  bool
)

func init() {
	pflag.StringVarP(&confile, "config", "c", "./config", "Configuration file")
	pflag.StringVarP(&logfile, "log", "l", "./server.log", "Logfile")
	pflag.BoolVarP(&useTLS, "tls", "t", false, "Use SSL")
	pflag.Parse()

	ch = InitLogger()

	// reading configuration file
	// see config.go
	ReadConfig()
}

func main() {
	// waits a second when server ends
	defer time.Sleep(time.Second)

	if useTLS {
		// configuring https server
		srvTLS = http.Server{
			Addr:      config.SSLPort,
			Handler:   &Handler{},
			TLSConfig: config.tlsConfig,
		}

		// configuring http2 support
		http2.ConfigureServer(&srvTLS, nil)
		go srvTLS.ListenAndServeTLS("", "")
		defer srvTLS.Close()
	}

	// http server
	srv = http.Server{
		Addr:    config.Port,
		Handler: &Handler{},
	}

	go srv.ListenAndServe()
	defer srv.Close()

	// signal handling (SIGINT)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	close(ch)
}
