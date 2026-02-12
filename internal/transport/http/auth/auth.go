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
		apperr.WriteError(w, err, h.logger)
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
		fmt.Println("something")
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	resp, err := h.authUsecase.VerifyUser(r.Context(), req.Email, req.Code, authUsecase.SessionMeta{
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
		Device:    meta.Device,
	})
	if err != nil {
		apperr.WriteError(w, err, h.logger)
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
		apperr.WriteError(w, err, h.logger)
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
		apperr.WriteError(w, err, h.logger)
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