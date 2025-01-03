package main

import (
	"log"
	"net/http"
	"time"
)

// Application
type application struct {
	config config
}

func (app *application) run(mux *http.ServeMux) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	return mux
}

// Config
type config struct {
	addr string
}
