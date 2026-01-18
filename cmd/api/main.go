package main

import (
	"log"
	"net/http"
	
	"github.com/All-Things-Muchiri/server/internal/config"
	"github.com/All-Things-Muchiri/server/internal/router"
	"github.com/joho/godotenv"
)

type application struct {
	config config.Config
}

func main() {
	// TODO: Create a util function to load and parse env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}
	
	cfg := config.LoadConfig()
	router := router.NewRouter()

	srv := &application{
		config: *cfg,
	}

	mux := router.Mount()
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
