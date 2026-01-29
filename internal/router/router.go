package router

import (
	"net/http"

	"github.com/All-Things-Muchiri/server/internal/auth"
	"github.com/All-Things-Muchiri/server/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppRouter struct {
	AuthProvider auth.AuthProvider
	WhHandler    *handler.WebhookHandler
}

func (a *AppRouter) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "status": "ok" }`))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Use(a.AuthProvider.RequireAuth())
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{ "status": "ok" }`))
		})
	})
	
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/clerk", a.WhHandler.CreateUser)
	})

	return r
}
