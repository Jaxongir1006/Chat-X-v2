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
			return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "user profile not found")
		}
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	media, err := u.userStore.GetProfileMedia(ctx, userID)
	if err != nil {
		u.logger.Warn().Err(err).Uint64("user_id", userID).Msg("failed to fetch profile media")
	}

	mediaDTOs := make([]UserProfileMediaDTO, 0, len(media))
	for _, m := range media {
		mediaDTOs = append(mediaDTOs, UserProfileMediaDTO{
			ID:           m.ID,
			MediaKey:     m.MediaKey,
			IsPrimary:    m.IsPrimary,
			DisplayOrder: m.DisplayOrder,
		})
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
			Media:        mediaDTOs,
			CreatedAt:    profile.CreatedAt,
			UpdatedAt:    profile.UpdatedAt,
		},
	}

	return &response, nil
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, userID uint64, req UpdateProfileRequest) error {
	return u.userStore.UpdateUserProfileFields(ctx, userID, req.FullName, req.Address, req.Bio)
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

func (u *UserUsecase) ChangePassword(ctx context.Context, userID uint64, req ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return apperr.New(apperr.CodeBadRequest, http.StatusBadRequest, "passwords do not match")
	}

	user, err := u.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if err := u.hasher.CheckPasswordHash(req.OldPassword, user.Password); err != nil {
		return apperr.New(apperr.CodeForbidden, http.StatusForbidden, "invalid old password")
	}

	newHash, err := u.hasher.Hash(req.NewPassword)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return u.userStore.UpdatePassword(ctx, userID, newHash)
}

func (u *UserUsecase) AddProfileMedia(ctx context.Context, userID uint64, req AddProfileMediaRequest) error {
	return u.userStore.AddProfileMedia(ctx, userID, req.MediaKey, req.IsPrimary)
}

func (u *UserUsecase) DeleteProfileMedia(ctx context.Context, userID uint64, mediaID uint64) error {
	return u.userStore.DeleteProfileMedia(ctx, userID, mediaID)
}

func (u *UserUsecase) SetPrimaryProfileMedia(ctx context.Context, userID uint64, mediaID uint64) error {
	return u.userStore.SetPrimaryProfileMedia(ctx, userID, mediaID)
}
