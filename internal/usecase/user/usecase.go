package user

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)

func (u *UserUsecase) GetMe(ctx context.Context, userID uint64) (*UserResponse, error) {
	user, err := u.userStore.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "user not found")
		}
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	profile, err := u.userStore.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "user profile not 	")
		}
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
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
			FullName:     profile.FullName,
			Address:      profile.Address,
			Bio:          profile.Bio,
			ProfileImage: profile.ProfileImage,
			CreatedAt:    profile.CreatedAt,
			UpdatedAt:    profile.UpdatedAt,
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

		err := sessTx.DeleteByUserID(ctx, userID)
		if err != nil {
			return err
		}

		err = userTx.DeleteUserProfile(ctx, userID)
		if err != nil {
			return err
		}

		err = userTx.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "failed to delete account", err)
	}

	return nil
}
