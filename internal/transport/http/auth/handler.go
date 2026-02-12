package auth

import (
	authUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/auth"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authUsecase *authUsecase.AuthUsecase
	logger      zerolog.Logger
}

func NewAuthHandler(authUsecase *authUsecase.AuthUsecase, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}
