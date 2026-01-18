package router

import (
	"net/http"

	"github.com/All-Things-Muchiri/server/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppRouter struct {
	authProvider auth.AuthProvider
}

func NewRouter(authProvider auth.AuthProvider) *AppRouter {
	return &AppRouter{
		authProvider: authProvider,
	}
}

func (a *AppRouter) Mount() http.Handler {
	r := chi.NewRouter()
	
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "status": "ok" }`))
	})
  
	r.Route("/auth", func(r chi.Router) {
		r.Use(a.authProvider.RequireAuth())
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{ "status": "ok" }`))
		})
	})

	return r
}