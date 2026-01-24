package router

import (
	"fmt"
	"net/http"

	"github.com/All-Things-Muchiri/server/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppRouter struct {
	authProvider auth.AuthProvider
}

func NewRouter(authProvider auth.AuthProvider) (*AppRouter, error) {
	if authProvider == nil {
		return nil, fmt.Errorf("authProvider is required but was nil")
	}
	return &AppRouter{
		authProvider: authProvider,
	}, nil
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
		r.Use(a.authProvider.RequireAuth())
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{ "status": "ok" }`))
		})
	})

	return r
}
