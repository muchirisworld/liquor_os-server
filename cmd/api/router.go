package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (a *application) mount() http.Handler {
	r := chi.NewRouter()
	
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
  
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "status": "ok" }`))
	})

	return r
}

func (a *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    a.config.addr,
		Handler: mux,
	}
	log.Printf("Server started on port %s", a.config.addr)
	
	return srv.ListenAndServe()
}
