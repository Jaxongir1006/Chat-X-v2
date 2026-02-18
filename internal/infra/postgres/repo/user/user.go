package userInfra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)

func (r *userRepo) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
	fmt.Println(userID)
	query := `SELECT username, phone, email, 
				verified, role, created_at, updated_at 
				FROM users WHERE id = $1`

	var result domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&result.Username,
		&result.Phone,
		&result.Email,
		&result.Verified,
		&result.Role,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "NOT FOUND")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	result.ID = userID
	return &result, nil
}

func (r *userRepo) GetUserProfileByUserID(ctx context.Context, userID uint64) (*domain.UserProfile, error) {
	query := `SELECT id, fullname, address, profile_image_link, created_at, updated_at FROM user_profile WHERE user_id = $1`

	var result domain.UserProfile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&result.ID, &result.FullName, &result.Address, &result.ProfileImage, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "NOT FOUND")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	result.UserID = userID
	return &result, nil
}

func (r *userRepo) UpdateUserProfile(ctx context.Context, userID uint64, profile *domain.UserProfile) error {
	query := `UPDATE user_profile SET fullname = $1, address = $2, profile_image_link = $3, updated_at = NOW() WHERE user_id = $4`

	_, err := r.db.ExecContext(ctx, query, profile.FullName, profile.Address, profile.ProfileImage, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *userRepo) DeleteUserProfile(ctx context.Context, userID uint64) error {
	query := `DELETE FROM user_profile WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}