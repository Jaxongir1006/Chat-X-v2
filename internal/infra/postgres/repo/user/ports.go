package userInfra

import (
	"context"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type UserStore interface {
	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	GetUserProfileByUserID(ctx context.Context, userID uint64) (*domain.UserProfile, error)
}
