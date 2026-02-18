package userUsecase

import (
	"context"
)

func (u *UserUsecase) GetMe(ctx context.Context, userID uint64) (*UserResponse, error) {
	user, err := u.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	profile, err := u.userStore.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Username:  user.Username,
		Verified:  user.Verified,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Profile: UserProfileResponse{
			ID: profile.ID,
			FullName: profile.FullName,
			Address: profile.Address,
			ProfileImage: profile.ProfileImage,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		},
	}

	return &response, nil
}
