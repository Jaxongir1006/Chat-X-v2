package user

import (
	userUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/user"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	usecase *userUsecase.UserUsecase
	logger  zerolog.Logger
}

func NewUserHandler(usecase *userUsecase.UserUsecase, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		usecase: usecase,
		logger:  logger,
	}
}
