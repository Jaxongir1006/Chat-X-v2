package media

import (
	"encoding/json"
	"net/http"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)


func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse the multipart form with a maximum memory of 5MB
    err := r.ParseMultipartForm(5 << 20)
    if err != nil {
        http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file from form: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    contentType := header.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    resp, err := h.usecase.UploadMedia(r.Context(), file, header.Filename, header.Size, contentType)
    if err != nil {
        apperr.WriteError(w, err, &h.logger)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(resp)
}

func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ObjectName string `json:"object_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteMedia(r.Context(), req.ObjectName); err != nil {
		apperr.WriteError(w, err, &h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Media deleted successfully"})
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}