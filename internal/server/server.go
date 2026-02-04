package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
)

type Server struct {
	mux  *http.ServeMux
	http *http.Server
}

func NewServer(cfg config.Server) *Server {
	mux := http.NewServeMux()

	s := &Server{
		mux: mux,
	}

	handler := middleware.Logging(mux)

	s.http = &http.Server{
		Addr:              cfg.Host + ":" + fmt.Sprint(cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

func (s *Server) Run() error {
	s.setupRoutes()
	log.Println("Server started on port 8080")
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.http.Shutdown(ctx)
}
