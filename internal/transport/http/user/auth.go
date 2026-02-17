package user

import (
	"encoding/json"
	"net/http"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
)

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	user, err := h.usecase.GetMe(r.Context(), userID)
	if err != nil {
		apperr.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}
