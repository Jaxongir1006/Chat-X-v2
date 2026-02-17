package server

import (
	"net/http"
)

func (s *Server) setupRoutes() {
	// health check in order to check if the server is running
	s.mux.HandleFunc("/health", healthCheck)

	// auth
	s.mux.HandleFunc("/api/v1/register", s.authHandler.Register)
	s.mux.HandleFunc("/api/v1/verify", s.authHandler.VerifyUser)
	s.mux.HandleFunc("/api/v1/login", s.authHandler.Login)
	s.mux.Handle("/api/v1/logout", s.authMiddleware.WrapAccess(http.HandlerFunc(s.authHandler.Logout)))
	s.mux.Handle("/api/v1/refresh", s.authMiddleware.WrapAccess(http.HandlerFunc(s.authHandler.Refresh)))

	// session
	s.mux.Handle("/api/v1/sessions", s.authMiddleware.WrapAccess(http.HandlerFunc(s.sessionHandler.Sessions)))
	s.mux.Handle("/api/v1/{session_id}/revoke", s.authMiddleware.WrapAccess(http.HandlerFunc(s.sessionHandler.RevokeSession)))

	// user
	s.mux.Handle("/api/v1/me", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.GetMe)))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
