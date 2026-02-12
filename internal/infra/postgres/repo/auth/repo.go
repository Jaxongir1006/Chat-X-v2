package authRepo

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/rs/zerolog"
)

type authRepo struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewAuthRepo(db *sql.DB, logger zerolog.Logger) *authRepo {
	return &authRepo{
		db:     db,
		logger: logger,
	}
}

func (r *authRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	query := `SELECT id, username, phone, email, verified, role,
				created_at, updated_at FROM users WHERE id = $1`

	var result domain.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
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

	return &result, nil
}

func (r *authRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, phone, verified, role,
				created_at, updated_at, password_hash FROM users WHERE email = $1`

	var result domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&result.ID,
		&result.Username,
		&result.Phone,
		&result.Verified,
		&result.Role,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Password,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "NOT FOUND")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	result.Email = email

	return &result, nil
}

func (r *authRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `SELECT id, username, email, verified, role,
				created_at, updated_at, password_hash FROM users WHERE phone = $1`

	var result domain.User

	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.Verified,
		&result.Role,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Password,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "NOT FOUND")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return &result, nil
}

func (r *authRepo) InsertUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, phone, email, password_hash, role, verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query, user.Username, user.Phone, user.Email, user.Password, user.Role, user.Verified)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *authRepo) DeleteUser(ctx context.Context, userID uint64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *authRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, phone, email, verified, role,
				created_at, updated_at, password_hash FROM users WHERE username = $1`

	var result domain.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&result.ID,
		&result.Username,
		&result.Phone,
		&result.Email,
		&result.Verified,
		&result.Role,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Password,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "NOT FOUND")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return &result, nil
}

func (r *authRepo) CreateUserProfile(ctx context.Context, userID uint64) error {
	query := `INSERT INTO profiles (user_id, created_at, updated_at) VALUES ($1, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *authRepo) GetUserProfile(ctx context.Context, userID uint64) (*domain.UserProfile, error) {
	query := `SELECT id, fullname, address, profile_image_link, created_at, updated_at FROM profiles WHERE user_id = $1`

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

func (r *authRepo) UpdateUserProfile(ctx context.Context, userID uint64, profile *domain.UserProfile) error {
	query := `UPDATE profiles SET fullname = $1, address = $2, profile_image_link = $3, updated_at = NOW() WHERE user_id = $4`

	_, err := r.db.ExecContext(ctx, query, profile.FullName, profile.Address, profile.ProfileImage, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *authRepo) DeleteUserProfile(ctx context.Context, userID uint64) error {
	query := `DELETE FROM profiles WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *authRepo) VerifyUser(ctx context.Context, email string) error {
	query := `UPDATE users SET verified = true WHERE email = $1`

	_, err := r.db.ExecContext(ctx, query, email)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}
