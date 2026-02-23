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
				verified, role, created_at, updated_at, password_hash 
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
		&result.Password,
	)
	if err != nil {
		return nil, err
	}
	result.ID = userID
	return &result, nil
}

func (r *userRepo) GetUserProfileByUserID(ctx context.Context, userID uint64) (*domain.UserProfile, error) {
	query := `SELECT fullname, address, bio, created_at, updated_at FROM user_profile WHERE user_id = $1`

	var result domain.UserProfile
	err := r.execer().QueryRowContext(ctx, query, userID).Scan(&result.FullName, &result.Address, &result.Bio, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	result.UserID = userID
	return &result, nil
}

func (r *userRepo) UpdateUserProfileFields(ctx context.Context, userID uint64, fullname, address, bio *string) error {
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

func (r *userRepo) UpdatePassword(ctx context.Context, userID uint64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.execer().ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (r *userRepo) AddProfileMedia(ctx context.Context, userID uint64, mediaKey string, isPrimary bool) error {
	query := `INSERT INTO user_profile_images (user_id, image_key, is_primary, display_order)
			  VALUES ($1, $2, $3, (SELECT COALESCE(MAX(display_order), 0) + 1 FROM user_profile_images WHERE user_id = $1))`
	_, err := r.execer().ExecContext(ctx, query, userID, mediaKey, isPrimary)
	return err
}

func (r *userRepo) GetProfileMedia(ctx context.Context, userID uint64) ([]domain.UserProfileMedia, error) {
	query := `SELECT id, user_id, image_key, is_primary, display_order, created_at 
			  FROM user_profile_images WHERE user_id = $1 ORDER BY display_order ASC`
	rows, err := r.execer().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []domain.UserProfileMedia
	for rows.Next() {
		var m domain.UserProfileMedia
		if err := rows.Scan(&m.ID, &m.UserID, &m.MediaKey, &m.IsPrimary, &m.DisplayOrder, &m.CreatedAt); err != nil {
			return nil, err
		}
		media = append(media, m)
	}
	return media, nil
}

func (r *userRepo) DeleteProfileMedia(ctx context.Context, userID uint64, mediaID uint64) error {
	query := `DELETE FROM user_profile_images WHERE id = $1 AND user_id = $2`
	_, err := r.execer().ExecContext(ctx, query, mediaID, userID)
	return err
}

func (r *userRepo) SetPrimaryProfileMedia(ctx context.Context, userID uint64, mediaID uint64) error {
	// First reset all to false, then set the specific one to true
	queryReset := `UPDATE user_profile_images SET is_primary = false WHERE user_id = $1`
	_, err := r.execer().ExecContext(ctx, queryReset, userID)
	if err != nil {
		return err
	}

	querySet := `UPDATE user_profile_images SET is_primary = true WHERE id = $1 AND user_id = $2`
	_, err = r.execer().ExecContext(ctx, querySet, mediaID, userID)
	return err
}
