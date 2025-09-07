package server

import (
	"L0/internal/config"
	"context"
	"net/http"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(config *config.HTTPServer, handler http.Handler) *Server {
	return &Server{
		HTTPServer: &http.Server{
			Addr:        config.Address,
			Handler:     handler,
			ReadTimeout: config.Timeout,
			IdleTimeout: config.Idle,
		},
	}
}

func (s *Server) Start() error {
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
