package server

import (
	"net/http"

)

func (s *Server) setupRoutes() {
	s.mux.Handle("/health", s.authMiddleware.AuthMiddleware(http.HandlerFunc(healthCheck)))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
