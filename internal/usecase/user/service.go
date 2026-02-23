package user

import (
	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	userInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/user"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/uow"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
	"github.com/rs/zerolog"
)

type UserUsecase struct {
	userStore userInfra.UserStore
	session   sessionInfra.SessionStore
	uow       uow.UnitOfWork
	hasher    security.Hasher
	logger    zerolog.Logger
}

func NewUserUsecase(userStore userInfra.UserStore, sessionStore sessionInfra.SessionStore, uow uow.UnitOfWork, hasher security.Hasher, logger zerolog.Logger) *UserUsecase {
	return &UserUsecase{
		userStore: userStore,
		session:   sessionStore,
		uow:       uow,
		hasher:    hasher,
		logger:    logger,
	}
}
