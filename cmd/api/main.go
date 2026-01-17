package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// TODO: Create a util function to load and parse env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}
	port := os.Getenv("PORT")

	cfg := &config{
		addr: port,
	}

	srv := &application{
		config: *cfg,
	}

	mux := srv.mount()
	log.Fatal(srv.run(mux))
}
