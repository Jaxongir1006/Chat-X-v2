package chat

import (
	chatRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/chat"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/uow"
	"github.com/rs/zerolog"
)

type ChatUsecase struct {
	chatStore chatRepo.ChatStore
	uow       uow.UnitOfWork
	logger    zerolog.Logger
}

func NewChatUsecase(chatStore chatRepo.ChatStore, uow uow.UnitOfWork, logger zerolog.Logger) *ChatUsecase {
	return &ChatUsecase{
		chatStore: chatStore,
		uow:       uow,
		logger:    logger,
	}
}
