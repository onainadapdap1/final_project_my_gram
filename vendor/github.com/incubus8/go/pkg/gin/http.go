package gin

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tylerb/graceful"
	"golang.org/x/net/http2"
)

// DefaultListenAddr represents the default address and port on which a server
// will listen, provided it is not overridden by setting the `ListenAddr` field
// on a `Config` struct.
const DefaultListenAddr = "0.0.0.0:8080"

// DefaultShutdownGracePeriod represents the default time in which the running
// process will allow outstanding http requests to complete before aborting
// them.  It will be used when a grace period of 0 is used, which normally
// signifies "no timeout" to our graceful shutdown package.  We choose not to
// provide a "no timeout" mode at present.  Feel free to set the value to a year
// or something if you need a timeout that is effectively "no timeout"; We
// believe that most servers should use a sane timeout and prefer one for the
// default configuration.
const DefaultShutdownGracePeriod = 10 * time.Second

// Config represents the configuration of an http server that can be provided to
// `Init`.
type Config struct {
	ShutdownGracePeriod time.Duration
	Handler             http.Handler
	ListenAddr          string
	TLSCert             string
	TLSKey              string
	OnStarting          func()
	OnStopping          func()
	OnStopped           func()
}

// Init starts an http server using the provided config struct.
//
// This method configures the process to listen for termination signals (SIGINT
// and SIGTERM) to trigger a graceful shutdown by way of the graceful package
// (https://github.com/tylerb/graceful).
func Run(conf Config) {
	srv := setup(conf)

	http2.ConfigureServer(srv.Server, nil)

	if conf.OnStarting != nil {
		conf.OnStarting()
	}

	var err error
	if conf.TLSCert != "" {
		err = srv.ListenAndServeTLS(conf.TLSCert, conf.TLSKey)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}

	if conf.OnStopped != nil {
		conf.OnStopped()
	}
	os.Exit(0)
}

func setup(conf Config) *graceful.Server {
	if conf.Handler == nil {
		panic("Handler must not be nil")
	}

	if conf.ListenAddr == "" {
		conf.ListenAddr = DefaultListenAddr
	}

	timeout := DefaultShutdownGracePeriod
	if conf.ShutdownGracePeriod != 0 {
		timeout = conf.ShutdownGracePeriod
	}

	return &graceful.Server{
		Timeout: timeout,

		Server: &http.Server{
			Addr:    conf.ListenAddr,
			Handler: conf.Handler,
		},

		ShutdownInitiated: func() {
			if conf.OnStopping != nil {
				conf.OnStopping()
			}
		},
	}
}
