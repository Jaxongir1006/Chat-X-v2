package media

import (
	usecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/media"
	"github.com/rs/zerolog"
)


type MediaHandler struct {
	usecase *usecase.MediaUsecase
	logger   zerolog.Logger
}

func NewMediaHandler(usecase *usecase.MediaUsecase, logger zerolog.Logger) *MediaHandler {
	return &MediaHandler{
		usecase: usecase,
		logger:  logger,
	}
}