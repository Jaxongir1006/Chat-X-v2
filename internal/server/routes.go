package server

import (
	"net/http"
)

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", healthCheck)
	s.mux.HandleFunc("/api/v1/register", s.authHandler.Register)
	s.mux.HandleFunc("/api/v1/verify", s.authHandler.VerifyUser)
	s.mux.HandleFunc("/api/v1/login", s.authHandler.Login)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
