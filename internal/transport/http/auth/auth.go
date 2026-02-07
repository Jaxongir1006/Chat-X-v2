package auth

import (
	"encoding/json"
	"net/http"

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


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("OK")
}
