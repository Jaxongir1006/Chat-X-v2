package sessionUsecase

import (
	"context"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type SessionUseCase interface {
	CreateSession(ctx context.Context, userId uint64, ip string, userAgent string, device string) (*domain.UserSession, error)
	ValidateAccess(ctx context.Context, accessToken string) (*domain.UserSession, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.UserSession, error)
	Logout(ctx context.Context, sessionID uint64) error
	LogoutAll(ctx context.Context, userID uint64) error
	EnforceMaxDevices(ctx context.Context, userID uint64, max int) error
}
