package server

import (
	"L0/internal/config"
	delivery "L0/internal/delivery/http"
	"context"
	"net/http"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(config *config.HTTPServer, handler *delivery.Handler) *Server {
	router := delivery.NewRouter(handler)

	return &Server{
		HTTPServer: &http.Server{
			Addr:        config.Address,
			Handler:     router,
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
