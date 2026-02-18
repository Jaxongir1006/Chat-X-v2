package sessionInfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type SessionStore interface {
	// used only when usecase needs transaction
	WithTx(tx *sql.Tx) *sessionRepo

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

	RotateRefresh(ctx context.Context, sessionID uint64, refresh string, refreshExp time.Time) error
	UpdateMeta(ctx context.Context, sessId uint64, device, ip, userAgent string, now time.Time) error

	RevokeByID(ctx context.Context, sessionID, userID uint64) error
	RevokeOthers(ctx context.Context, userID uint64, currentSessionID uint64) error
	RevokeAllByUserID(ctx context.Context, userID uint64) error
	RevokeAllExceptCurrent(ctx context.Context, userID uint64, sessionID uint64) error
}
