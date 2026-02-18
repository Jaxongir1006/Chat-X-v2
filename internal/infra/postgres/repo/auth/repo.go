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
	tx     *sql.Tx
	logger zerolog.Logger
}

func NewAuthRepo(db *sql.DB, logger zerolog.Logger) *authRepo {
	return &authRepo{db: db, logger: logger}
}

func (r *authRepo) WithTx(tx *sql.Tx) *authRepo {
	return &authRepo{db: r.db, tx: tx, logger: r.logger}
}

func (r *authRepo) execer() interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}


func (r *authRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	query := `SELECT id, username, phone, email, verified, role,
				created_at, updated_at FROM users WHERE id = $1`

	var result domain.User

	err := r.execer().QueryRowContext(ctx, query, id).Scan(
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
	err := r.execer().QueryRowContext(ctx, query, email).Scan(
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

	err := r.execer().QueryRowContext(ctx, query, phone).Scan(
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
	query := `INSERT INTO users (username, phone, email, password_hash, role, verified, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.execer().ExecContext(ctx, query, user.Username, user.Phone, user.Email, user.Password, user.Role, user.Verified)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (r *authRepo) DeleteUser(ctx context.Context, userID uint64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.execer().ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (r *authRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, phone, email, verified, role,
				created_at, updated_at, password_hash FROM users WHERE username = $1`

	var result domain.User

	err := r.execer().QueryRowContext(ctx, query, username).Scan(
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
	query := `INSERT INTO user_profile (user_id, created_at, updated_at) VALUES ($1, NOW(), NOW())`

	_, err := r.execer().ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (r *authRepo) VerifyUser(ctx context.Context, email string) error {
	query := `UPDATE users SET verified = true WHERE email = $1`

	_, err := r.execer().ExecContext(ctx, query, email)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (r *authRepo) RestartUnverified(ctx context.Context, id uint64, username, phone, hashed string) error {
	query := `UPDATE users SET username = $2, phone = $3, password_hash = $4 WHERE id = $1`

	_, err := r.execer().ExecContext(ctx, query, id, username, phone, hashed)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}