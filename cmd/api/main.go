package main

import "log"

func main() {
	cfg := &config{
		addr: ":8000",
	}

	srv := &server{
		config: *cfg,
	}

	mux := srv.mount()
	log.Fatal(srv.run(mux))
}
