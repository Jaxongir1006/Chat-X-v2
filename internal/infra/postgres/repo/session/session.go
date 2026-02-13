package sessionInfra

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)

func (r *sessionRepo) GetAllValidSessionsByUserId(ctx context.Context, userID uint64) ([]domain.UserSession, error) {
	query := `SELECT id, refresh_token, refresh_token_expires_at, access_token, access_token_expires_at, 
				last_used_at, ip_address, user_agent, device, created_at, updated_at
				FROM serefresh_token_expires_atssions WHERE user_id = $1 AND  > NOW()`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Fatal().Err(err).Msg("could not close rows")
		}
	}()

	sessions := make([]domain.UserSession, 0)

	for rows.Next() {
		var s domain.UserSession
		if err := rows.Scan(&s.ID, &s.RefreshToken, &s.RefreshTokenExp, &s.AccessToken,
			&s.AccessTokenExp, &s.LastUsedAt, &s.IPAddress, &s.UserAgent,
			&s.Device, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
		s.UserID = userID
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return sessions, nil
}

func (r *sessionRepo) GetByAccessToken(ctx context.Context, accessToken string) (*domain.UserSession, error) {
	query := `SELECT id, user_id, refresh_token, refresh_token_expires_at, access_token, access_token_expires_at, 
				last_used_at, ip_address, user_agent, device, created_at, updated_at
				FROM sessions WHERE access_token = $1 AND access_token_expires_at > NOW()`

	var result domain.UserSession

	err := r.db.QueryRowContext(ctx, query, accessToken).Scan(&result.ID, &result.UserID, &result.RefreshToken, &result.RefreshTokenExp, &result.AccessToken,
		&result.AccessTokenExp, &result.LastUsedAt, &result.IPAddress, &result.UserAgent,
		&result.Device, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return &result, nil
}

func (r *sessionRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.UserSession, error) {
	query := `SELECT id, user_id, refresh_token, refresh_token_expires_at, access_token, access_token_expires_at, 
				last_used_at, ip_address, user_agent, device, created_at, updated_at
				FROM sessions WHERE refresh_token = $1 AND refresh_token_expires_at > NOW()`

	var result domain.UserSession

	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&result.ID, &result.UserID, &result.RefreshToken, &result.RefreshTokenExp, &result.AccessToken,
		&result.AccessTokenExp, &result.LastUsedAt, &result.IPAddress, &result.UserAgent,
		&result.Device, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED")
	}
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return &result, nil
}

func (r *sessionRepo) Create(ctx context.Context, s *domain.UserSession) error {
	query := `INSERT INTO sessions (user_id, refresh_token, refresh_token_expires_at, access_token, access_token_expires_at,
            	last_used_at, ip_address, user_agent, device)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
				RETURNING id, created_at, updated_at
				`

	_, err := r.db.ExecContext(ctx, query, s.UserID, s.RefreshToken, s.RefreshTokenExp, s.AccessToken,
		s.AccessTokenExp, s.LastUsedAt, s.IPAddress, s.UserAgent, s.Device)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) UpdateTokens(ctx context.Context, sessionID uint64, access string, accessExp time.Time, refresh string, refreshExp time.Time) error {
	query := `UPDATE sessions SET access_token = $1, access_token_expires_at = $2, refresh_token = $3, refresh_token_expires_at = $4, updated_at = NOW()
				WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query, access, accessExp, refresh, refreshExp, sessionID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) DeleteByID(ctx context.Context, sessionID uint64) error {
	query := `DELETE FROM sessions WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) DeleteByUserID(ctx context.Context, userID uint64) error {
	query := `DELETE FROM sessions WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) DeleteOldestValidSession(ctx context.Context, userID uint64) error {
	query := `DELETE FROM sessions
				WHERE id = (
					SELECT id
					FROM sessions
					WHERE user_id = $1 AND refresh_token_expires_at > NOW()
					ORDER BY created_at
					LIMIT 1
				)`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) DeleteExpiredRefreshSessionsByUserID(ctx context.Context, userID uint64) error {
	query := `DELETE FROM sessions WHERE user_id = $1 AND refresh_token_expires_at < NOW()`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) RotateRefresh(ctx context.Context, sessionID uint64, refresh string, refreshExp time.Time) error {
	query := `UPDATE sessions SET refresh_token = $1, refresh_token_expires_at = $2, updated_at = NOW()
				WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, refresh, refreshExp, sessionID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}

func (r *sessionRepo) UpdateMeta(ctx context.Context, sessId int, device, ip, userAgent string, now time.Time) error {
	query := `UPDATE sessions SET device = $1, ip_address = $2, user_agent = $3, updated_at = $4
				WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query, device, ip, userAgent, now, sessId)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return nil
}