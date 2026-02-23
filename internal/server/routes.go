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
	s.mux.Handle("/api/v1/me/profile", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.UpdateProfile)))
	s.mux.Handle("/api/v1/me/delete", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.DeleteAccount)))
	s.mux.Handle("/api/v1/me/password", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.ChangePassword)))
	s.mux.Handle("/api/v1/me/profile/media", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.AddProfileMedia)))
	s.mux.Handle("/api/v1/me/profile/media/delete", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.DeleteProfileMedia)))
	s.mux.Handle("/api/v1/me/profile/media/primary", s.authMiddleware.WrapAccess(http.HandlerFunc(s.userHandler.SetPrimaryProfileMedia)))

	// chat
	s.mux.Handle("/api/v1/chat/conversations", s.authMiddleware.WrapAccess(http.HandlerFunc(s.chatHandler.GetConversations)))
	s.mux.Handle("/api/v1/chat/dm", s.authMiddleware.WrapAccess(http.HandlerFunc(s.chatHandler.StartDM)))
	s.mux.Handle("/api/v1/chat/group", s.authMiddleware.WrapAccess(http.HandlerFunc(s.chatHandler.CreateGroup)))
	s.mux.Handle("/api/v1/chat/messages/history", s.authMiddleware.WrapAccess(http.HandlerFunc(s.chatHandler.GetMessages)))

	// media 
	s.mux.Handle("/api/v1/media/upload", s.authMiddleware.WrapAccess(http.HandlerFunc(s.mediaHandler.UploadMedia)))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
