package authUsecase

import (
	authRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/auth"
	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	redisStore "github.com/Jaxongir1006/Chat-X-v2/internal/infra/redis/store"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
)

type AuthUsecase struct {
	authStore authRepo.AuthStore
	session   sessionInfra.SessionStore
	redis     redisStore.OTPStore
	hasher    security.Hasher
	token     security.TokenStore
}

func NewAuthUsecase(authStore authRepo.AuthStore, session sessionInfra.SessionStore, redis redisStore.OTPStore, token security.TokenStore, hasher security.Hasher) *AuthUsecase {
	return &AuthUsecase{
		authStore: authStore,
		session:   session,
		redis:     redis,
		token:     token,
		hasher:    hasher,
	}
}
