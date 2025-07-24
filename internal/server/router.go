package server

import "net/http"

// newRouter sets up and returns the *http.ServeMux with all routes.
func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// health check / root
	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /healthz", handleHealth)

	// TODO: add more routes here without touching server.go
	return mux
}
