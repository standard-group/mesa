package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	http *http.Server
}

// New constructs and returns the HTTP server
func New(port int) *Server {
	mux := newRouter() // defined in router.go

	return &Server{
		http: &http.Server{
			Addr:         ":" + strconv.Itoa(port),
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() error {
	log.Info().Str("addr", s.http.Addr).Msg("HTTP server starting")
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) {
	const timeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	log.Info().Msg("Shutting down HTTP server")
	if err := s.http.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}
}
