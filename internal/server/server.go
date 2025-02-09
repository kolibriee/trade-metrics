package server

import (
	"context"
	"net/http"

	"github.com/kolibriee/trade-metrics/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Server, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
