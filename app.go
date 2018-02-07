package atc

import (
	"net/http"
	"time"

	"github.com/adolphlxm/atc/logs"
)

var HttpAPP *App

// A APP defines parameters for running an ATC server.
type App struct {
	Handler *HandlerRouter
	// A Server defines parameters for running an HTTP server.
	// The zero value for Server is a valid configuration.
	Server *http.Server
}

// NewApp returns a new atc application.
func NewApp() *App {
	h, _ := NewHandlerRouter()
	app := &App{Handler: h, Server: &http.Server{}}
	return app
}

// Run atc application.
func (a *App) Run() {
	// Defines parameters for running an HTTP server.
	addr := Aconfig.HTTPAddr + ":" + Aconfig.HTTPPort
	a.Server = &http.Server{
		Addr:         addr,
		Handler:      a.Handler,
		ReadTimeout:  time.Duration(Aconfig.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(Aconfig.HTTPWriteTimeout) * time.Second,
	}

	//appRun := make(chan bool, 1)
	//TODO https server
	end := make(chan struct{})
	//http server
	go func() {
		logs.Trace("http: starting...")
		close(end)
		// ListenAndServe listens on the TCP network address srv.Addr and then
		// calls Serve to handle requests on incoming connections.
		err := a.Server.ListenAndServe()
		if err != nil {
			logs.Errorf("http: ListenAndServe err: %v", err)
		} else {
			logs.Tracef("http: Running on %s", addr)
		}
	}()

	<-end
}
