package authRepo

import (
	"context"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type AuthStore interface {
	GetByID(ctx context.Context, id uint64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	InsertUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID uint64) error
	CreateUserProfile(ctx context.Context, userID uint64) error
	GetUserProfile(ctx context.Context, userID uint64) (*domain.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID uint64, profile *domain.UserProfile) error
	DeleteUserProfile(ctx context.Context, userID uint64) error
	VerifyUser(ctx context.Context, email string) error
}
