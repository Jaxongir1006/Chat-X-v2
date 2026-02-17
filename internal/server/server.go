package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/auth"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/session"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/user"
	"github.com/rs/zerolog"
)

type Server struct {
	mux            *http.ServeMux
	http           *http.Server
	authMiddleware *middleware.AuthMiddleware
	authHandler    *auth.AuthHandler
	sessionHandler *session.SessionHandler
	userHandler    *user.UserHandler
	logger         zerolog.Logger
}

func NewServer(cfg config.Server, authMiddleware *middleware.AuthMiddleware, logger zerolog.Logger,
	authHandler *auth.AuthHandler, sessionHandler *session.SessionHandler, userHandler *user.UserHandler) *Server {
	mux := http.NewServeMux()

	s := &Server{
		mux:            mux,
		authMiddleware: authMiddleware,
		authHandler:    authHandler,
		sessionHandler: sessionHandler,
		userHandler:    userHandler,
		logger:         logger,
	}

	var handler http.Handler = mux
	handler = middleware.MetaMiddleware(handler)
	handler = middleware.Logging(logger, handler)

	s.http = &http.Server{
		Addr:              cfg.Host + ":" + fmt.Sprint(cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

func (s *Server) Run() error {
	s.setupRoutes()
	s.logger.Info().Msg("Starting the HTTP server...")
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down the HTTP server...")
	return s.http.Shutdown(ctx)
}
