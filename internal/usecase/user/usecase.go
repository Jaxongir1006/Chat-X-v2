package userUsecase

import (
	"context"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

func (u *UserUsecase) GetMe(ctx context.Context, userID uint64) (*domain.User, error) {
	return u.userStore.GetUserByID(ctx, userID)
}
