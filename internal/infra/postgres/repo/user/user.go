package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

func (r *userRepo) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
	fmt.Println(userID)
	query := `SELECT username, phone, email, 
				verified, role, created_at, updated_at 
				FROM users WHERE id = $1`

	var result domain.User
	err := r.execer().QueryRowContext(ctx, query, userID).Scan(
		&result.Username,
		&result.Phone,
		&result.Email,
		&result.Verified,
		&result.Role,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	result.ID = userID
	return &result, nil
}

func (r *userRepo) GetUserProfileByUserID(ctx context.Context, userID uint64) (*domain.UserProfile, error) {
	query := `SELECT fullname, address, bio, profile_image_key, created_at, updated_at FROM user_profile WHERE user_id = $1`

	var result domain.UserProfile
	err := r.execer().QueryRowContext(ctx, query, userID).Scan(&result.FullName, &result.Address, &result.Bio, &result.ProfileImage, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	result.UserID = userID
	return &result, nil
}

func (r *userRepo) UpdateUserProfileFields(ctx context.Context, userID uint64, fullname, address, profileImage, bio *string) error {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if fullname != nil {
		setClauses = append(setClauses, fmt.Sprintf("fullname = $%d", idx))
		args = append(args, *fullname)
		idx++
	}

	if address != nil {
		setClauses = append(setClauses, fmt.Sprintf("address = $%d", idx))
		args = append(args, *address)
		idx++
	}

	if bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", idx))
		args = append(args, *bio)
		idx++
	}

	if profileImage != nil {
		setClauses = append(setClauses, fmt.Sprintf("profile_image_key = $%d", idx))
		args = append(args, *profileImage)
		idx++
	}

	if len(setClauses) == 1 {
		return nil
	}

	where := fmt.Sprintf("user_id = $%d", idx)
	args = append(args, userID)

	query := `
		UPDATE user_profile
		SET ` + strings.Join(setClauses, ", ") + `
		WHERE ` + where

	_, err := r.execer().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) DeleteUserProfile(ctx context.Context, userID uint64) error {
	query := `DELETE FROM user_profile WHERE user_id = $1`

	_, err := r.execer().ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, userID uint64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.execer().ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}
