package user

import (
	"context"
	"database/sql"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type UserStore interface {
	// used only when usecase needs transaction
	WithTx(tx *sql.Tx) *userRepo

	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	GetUserProfileByUserID(ctx context.Context, userID uint64) (*domain.UserProfile, error)
	UpdateUserProfileFields(ctx context.Context, userID uint64, fullname, address, bio *string) error
	DeleteUser(ctx context.Context, userID uint64) error
	DeleteUserProfile(ctx context.Context, userID uint64) error

	// Password and Media
	UpdatePassword(ctx context.Context, userID uint64, passwordHash string) error
	AddProfileMedia(ctx context.Context, userID uint64, mediaKey string, isPrimary bool) error
	GetProfileMedia(ctx context.Context, userID uint64) ([]domain.UserProfileMedia, error)
	DeleteProfileMedia(ctx context.Context, userID uint64, mediaID uint64) error
	SetPrimaryProfileMedia(ctx context.Context, userID uint64, mediaID uint64) error
}
