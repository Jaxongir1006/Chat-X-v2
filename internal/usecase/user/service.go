package userUsecase

import (
	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/user"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/uow"
	"github.com/rs/zerolog"
)

type UserUsecase struct {
	userStore userInfra.UserStore
	session   sessionInfra.SessionStore
	uow       uow.UnitOfWork
	logger    zerolog.Logger
}

func NewUserUsecase(userStore userInfra.UserStore, sessionStore sessionInfra.SessionStore, uow uow.UnitOfWork, logger zerolog.Logger) *UserUsecase {
	return &UserUsecase{
		userStore: userStore,
		session:   sessionStore,
		uow:       uow,
		logger:    logger,
	}
}
