package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
	authUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/auth"
)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authUsecase.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: Invalid Json", http.StatusBadRequest)
		return
	}

	if err := h.authUsecase.Register(r.Context(), req); err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.logger.Info().Msg("New user registered successfully")
	err := json.NewEncoder(w).Encode(map[string]any{
		"message": "User registered successfully",
		"success": true,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authUsecase.VerifyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: Invalid Json", http.StatusBadRequest)
		return
	}

	meta, ok := middleware.MetaFromContext(r.Context())
	if !ok {
		h.logger.Error().Msg("Failed to get meta from context")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	resp, err := h.authUsecase.VerifyUser(r.Context(), req.Email, req.Code, authUsecase.SessionMeta{
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
		Device:    meta.Device,
	})
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.logger.Info().Msg(fmt.Sprintf("User with the email %s verified successfully", req.Email))
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authUsecase.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: Invalid Json", http.StatusBadRequest)
		return
	}

	meta, ok := middleware.MetaFromContext(r.Context())
	if !ok {
		h.logger.Error().Msg("Failed to get meta from context")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	resp, err := h.authUsecase.Login(r.Context(), req, authUsecase.SessionMeta{
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
		Device:    meta.Device,
	})
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.logger.Info().Msg("User logged in successfully")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authUsecase.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: Invalid Json", http.StatusBadRequest)
		return
	}

	meta, ok := middleware.MetaFromContext(r.Context())
	if !ok {
		h.logger.Error().Msg("Failed to get meta from context")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	resp, err := h.authUsecase.Refresh(r.Context(), req, authUsecase.SessionMeta{
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
		Device:    meta.Device,
	})
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.logger.Info().Msg("User logged in successfully")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()

	var req authUsecase.LogoutRequest
	if err := dec.Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Bad request: Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	op := req.Operation

	var err error
	switch op {
	case "all":
		err = h.authUsecase.LogoutAll(ctx, userID)

	case "one":
		if req.SessionID == nil {
			http.Error(w, "Bad request: Missing session_id", http.StatusBadRequest)
			return
		}
		err = h.authUsecase.LogoutFromCurrent(ctx, *req.SessionID, userID)

	case "except-current":
		if req.SessionID == nil {
			h.logger.Error().Msg("Invalid logout operation and missing session_id")
			http.Error(w, "Bad request: Missing session_id", http.StatusBadRequest)
			return
		}
		err = h.authUsecase.LogOutAllExceptCurrent(ctx, *req.SessionID, userID)

	default:
		if req.SessionID == nil {
			h.logger.Error().Msg("Invalid logout operation and missing session_id")
			http.Error(w, "Bad Request: Missing session_id", http.StatusBadRequest)
			return
		}
		err = h.authUsecase.LogoutFromCurrent(ctx, userID, *req.SessionID)
	}

	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "User logged out successfully",
		"success": true,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode response")
		return
	}

	h.logger.Info().Uint64("user_id", userID).Str("operation", op).Msg("User logged out")
}
