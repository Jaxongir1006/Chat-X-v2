package media

import (
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/minio"
	"github.com/rs/zerolog"
)


type MediaUsecase struct {
	storage minio.ObjectStorage
	logger  zerolog.Logger
}

func NewMediaUsecase(storage minio.ObjectStorage, logger zerolog.Logger) *MediaUsecase {
	return &MediaUsecase{
		storage: storage,
		logger:  logger,
	}
}