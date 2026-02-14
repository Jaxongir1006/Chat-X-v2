package sessionUsecase

import (
	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
)

type SessionUsecase struct {
	sessionStore sessionInfra.SessionStore
	token        *security.Token
	maxDevices   int
}

func NewSessionService(store sessionInfra.SessionStore, token *security.Token, maxDevices int) *SessionUsecase {
	return &SessionUsecase{
		sessionStore: store,
		token:        token,
		maxDevices:   maxDevices,
	}
}
