package session

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
)

func (h *SessionHandler) Sessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	userSessions, err := h.sessionUsecase.GetSessionsByUserID(r.Context(), userID)
	if err != nil {
		apperr.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userSessions)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}

func (h *SessionHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	key := r.URL.Query().Get("session_id")
	if key == "" {
		http.Error(w, "Bad request: Missing session_id", http.StatusBadRequest)
		return
	}

	sessID, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		http.Error(w, "Bad request: Invalid session_id", http.StatusBadRequest)
		return
	}

	if err := h.sessionUsecase.RevokeSessionByID(r.Context(), sessID, userID); err != nil {
		apperr.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.logger.Info().Msg("Session revoked successfully")
	err = json.NewEncoder(w).Encode(map[string]any{
		"message": "Session revoked successfully",
		"success": true,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}
