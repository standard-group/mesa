package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/standard-group/mesa/internal/auth"
)

func NewServer() *http.Server {
	r := chi.NewRouter()
	r.Use(auth.AuthMiddleware)
	r.Post("/api/v1/register", RegisterHandler)
	r.Post("/api/v1/login", LoginHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	return srv
}
