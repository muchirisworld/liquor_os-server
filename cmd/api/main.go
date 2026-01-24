package main

import (
	"log"
	"net/http"
	"os"

	"github.com/All-Things-Muchiri/server/internal/auth"
	"github.com/All-Things-Muchiri/server/internal/clerk"
	"github.com/All-Things-Muchiri/server/internal/config"
	"github.com/All-Things-Muchiri/server/internal/router"
	"github.com/joho/godotenv"
)

type application struct {
	config       config.Config
	authProvider auth.AuthProvider
}

func main() {
	// TODO: Create a util function to load and parse env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}

	secret := os.Getenv("AUTH_SECRET_KEY")
	if secret == "" {
		log.Fatalf("AUTH_SECRET_KEY environment variable is required but not set")
	}

	authConfig := &config.AuthConfig{
		SecretKey: secret,
	}

	cfg := config.LoadConfig(authConfig)
	clerkProvider := clerk.NewProvider(*cfg.AuthConfig)
	apiRouter, err := router.NewRouter(clerkProvider)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	srv := &application{
		config: *cfg,
	}

	mux := apiRouter.Mount()
	log.Fatal(srv.run(mux))
}

func (a *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    a.config.Addr,
		Handler: mux,
	}
	log.Printf("Server started on port %s", a.config.Addr)

	return srv.ListenAndServe()
}
