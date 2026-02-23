package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
	chatUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/chat"
)

func (h *ChatHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	convs, err := h.usecase.GetConversations(r.Context(), userID)
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}

func (h *ChatHandler) StartDM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	var req chatUsecase.StartDMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	resp, err := h.usecase.StartDM(r.Context(), userID, req.UserID)
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ChatHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	var req chatUsecase.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	resp, err := h.usecase.CreateGroup(r.Context(), userID, req)
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	convIDStr := r.URL.Query().Get("conversation_id")
	convID, err := strconv.ParseUint(convIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	messages, err := h.usecase.GetMessages(r.Context(), userID, convID, limit, offset)
	if err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
