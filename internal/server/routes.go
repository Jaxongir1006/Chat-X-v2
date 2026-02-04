package server

import "net/http"

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", healthCheck)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
