package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server wraps HTTP server with configuration
type Server struct {
	httpServer *http.Server
	port       int
}

// NewServer creates a new HTTP server
func NewServer(port int, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		port: port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("Starting HTTP server on port %d\n", s.port)
	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

// Port returns the server port
func (s *Server) Port() int {
	return s.port
}
