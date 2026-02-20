package auth

import (
	authRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/auth"
	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/uow"
	redisStore "github.com/Jaxongir1006/Chat-X-v2/internal/infra/redis/store"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
	"github.com/rs/zerolog"
)

type AuthUsecase struct {
	uow        uow.UnitOfWork
	authStore  authRepo.AuthStore
	session    sessionInfra.SessionStore
	redis      redisStore.OTPStore
	hasher     security.Hasher
	token      security.TokenStore
	logger     zerolog.Logger
	codeHasher security.CodeHasher
}

func NewAuthUsecase(authStore authRepo.AuthStore,
	session sessionInfra.SessionStore, redis redisStore.OTPStore,
	token security.TokenStore, hasher security.Hasher,
	logger zerolog.Logger, codeHasher security.CodeHasher, uow uow.UnitOfWork) *AuthUsecase {

	return &AuthUsecase{
		authStore:  authStore,
		session:    session,
		redis:      redis,
		token:      token,
		hasher:     hasher,
		logger:     logger,
		codeHasher: codeHasher,
		uow:        uow,
	}
}
