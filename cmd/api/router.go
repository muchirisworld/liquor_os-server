package main

import "net/http"

type server struct {
	config config
}

type config struct {
	addr string
}

func (s *server) mount() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "status": "ok" }`))
	})

	return r
}

func (s *server) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    s.config.addr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
