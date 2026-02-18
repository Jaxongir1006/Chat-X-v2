package authRepo

import (
	"context"
	"database/sql"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type AuthStore interface {
	// used only when usecase needs transaction
	WithTx(tx *sql.Tx) *authRepo

	GetByID(ctx context.Context, id uint64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	InsertUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID uint64) error
	CreateUserProfile(ctx context.Context, userID uint64) error
	VerifyUser(ctx context.Context, email string) error
	RestartUnverified(ctx context.Context, id uint64, username, phone, hashed string) error
}
