package chat

import (
	chatUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/chat"
	"github.com/rs/zerolog"
)

type ChatHandler struct {
	usecase *chatUsecase.ChatUsecase
	logger  zerolog.Logger
}

func NewChatHandler(usecase *chatUsecase.ChatUsecase, logger zerolog.Logger) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
		logger:  logger,
	}
}
