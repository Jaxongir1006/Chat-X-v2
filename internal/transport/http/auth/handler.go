package auth

import (
	authUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/auth"
)

type AuthHandler struct {
	authUsecase *authUsecase.AuthUsecase
}

func NewAuthHandler(authUsecase *authUsecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}
