package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)


func (h *MediaUsecase) UploadMedia(ctx context.Context, file multipart.File, filename string, size int64, contentType string) (*UploadMediaResponse, error) {
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	_, err := h.storage.Upload(ctx, objectName, file, size, contentType)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	mediaURL, err := h.storage.PresignGet(ctx, objectName, 24*time.Hour)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return &UploadMediaResponse{
		MediaURL: mediaURL,
		ObjectName: objectName,
	}, nil
}


func (h *MediaUsecase) DeleteMedia(ctx context.Context, objectName string) error {
	if objectName == "" {
		return apperr.New(apperr.CodeBadRequest, http.StatusBadRequest, "object name is required")
	}

	if err := h.storage.Delete(ctx, objectName); err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}