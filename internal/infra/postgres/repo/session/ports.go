package sessionInfra

import (
	"context"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type SessionStore interface {
	// lists sessions where refresh is still valid (or revoked_at is null)
	GetAllValidSessionsByUserId(ctx context.Context, userID uint64) ([]domain.UserSession, error)

	// token lookups (index tokens in DB)
	GetByAccessToken(ctx context.Context, accessToken string) (*domain.UserSession, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.UserSession, error)

	// CRUD
	Create(ctx context.Context, s *domain.UserSession) error
	UpdateTokens(ctx context.Context, sessionID uint64, access string, accessExp time.Time, refresh string, refreshExp time.Time) error

	DeleteByID(ctx context.Context, sessionID uint64) error
	DeleteByUserID(ctx context.Context, userID uint64) error

	// device/session limit helpers
	DeleteOldestValidSession(ctx context.Context, userID uint64) error

	// cleanup
	DeleteExpiredRefreshSessionsByUserID(ctx context.Context, userID uint64) error
}
