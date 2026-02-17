package userInfra

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)

func (r *userRepo) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
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

	return nil, nil
}
