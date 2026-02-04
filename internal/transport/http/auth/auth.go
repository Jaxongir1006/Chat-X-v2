package auth

import (
	"encoding/json"
	"net/http"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// var req RegisterRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, "Bad request: Invalid Json", http.StatusBadRequest)
	// 	return
	// }

	// resp, err := h.Svc.Register(req)
	// if err != nil {
	// 	apperr.WriteError(w, err)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("OK")
}
