package userUsecase

import (
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/user"
	"github.com/rs/zerolog"
)

type UserUsecase struct {
	userStore userInfra.UserStore
	logger    zerolog.Logger
}

func NewUserUsecase(userStore userInfra.UserStore, logger zerolog.Logger) *UserUsecase {
	return &UserUsecase{
		userStore: userStore,
		logger:    logger,
	}
}
