package userUsecase

import (
	"context"
	"database/sql"
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
			FullName: profile.FullName,
			Address: profile.Address,
			Bio: profile.Bio,
			ProfileImage: profile.ProfileImage,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		},
	}

	return &response, nil
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, userID uint64, req UpdateProfileRequest) error {
	return u.userStore.UpdateUserProfileFields(ctx, userID, req.FullName, req.Address, req.ProfileImage, req.Bio)
}


func (u *UserUsecase) DeleteAccount(ctx context.Context, userID uint64) error {
	err := u.uow.Do(ctx, func(tx *sql.Tx) error {
		userTx := u.userStore.WithTx(tx)
		sessTx := u.session.WithTx(tx)

		return nil
	})

	
	return nil
}