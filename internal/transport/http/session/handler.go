package session

import (
	sessionUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/session"
	"github.com/rs/zerolog"
)

type SessionHandler struct {
	sessionUsecase *sessionUsecase.SessionUsecase
	logger         zerolog.Logger
}

func NewSessionHandler(authUsecase *sessionUsecase.SessionUsecase, logger zerolog.Logger) *SessionHandler {
	return &SessionHandler{
		sessionUsecase: authUsecase,
		logger:         logger,
	}
}
