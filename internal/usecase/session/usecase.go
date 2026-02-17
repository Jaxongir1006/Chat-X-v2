package sessionUsecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)

func (s *SessionUsecase) CreateSession(
	ctx context.Context,
	userId uint64,
	ip, userAgent, device string,
) (*domain.UserSession, error) {
	sessions, err := s.sessionStore.GetAllValidSessionsByUserId(ctx, userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		sessions = []domain.UserSession{}
	}

	max := s.maxDevices
	if max <= 0 {
		max = 5
	}

	// enforce max devices: if already 5 => delete oldest, then create new
	if len(sessions) >= max {
		if err := s.sessionStore.DeleteOldestValidSession(ctx, userId); err != nil {
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
	}

	access, accessExp, err := s.token.GenerateAccessToken(strconv.FormatUint(userId, 10))
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	refresh, refreshExp, err := s.token.GenerateRefreshToken(strconv.FormatUint(userId, 10))
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	now := time.Now()

	newSess := &domain.UserSession{
		UserID:          userId,
		AccessToken:     access,
		AccessTokenExp:  accessExp,
		RefreshToken:    refresh,
		RefreshTokenExp: refreshExp,
		IPAddress:       ip,
		UserAgent:       userAgent,
		Device:          device,
		LastUsedAt:      &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.sessionStore.Create(ctx, newSess); err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return newSess, nil
}

func (s *SessionUsecase) ValidateAccess(ctx context.Context, accessToken string) (*domain.UserSession, error) {
	claims, err := s.token.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", err)
	}

	sess, err := s.sessionStore.GetByAccessToken(ctx, accessToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", err)
		}
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// extra safety: match userID from token with DB
	dbUserID := strconv.FormatUint(sess.UserID, 10)
	if claims.UserID != dbUserID {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED")
	}

	// also optional: check DB expiries (in case you revoked/changed logic)
	if time.Now().After(sess.AccessTokenExp) {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "ACCESS TOKEN EXPIRED")
	}

	return sess, nil
}

func (s *SessionUsecase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.UserSession, error) {
	claims, err := s.token.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", err)
	}

	sess, err := s.sessionStore.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", err)
		}
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	dbUserID := strconv.FormatUint(sess.UserID, 10)
	if claims.UserID != dbUserID {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED")
	}

	now := time.Now()
	if now.After(sess.RefreshTokenExp) {
		// refresh expired -> session dead
		_ = s.sessionStore.DeleteByID(ctx, sess.ID)
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "REFRESH TOKEN EXPIRED")
	}

	// new access token always
	newAccess, newAccessExp, err := s.token.GenerateAccessToken(dbUserID)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// rotate refresh token (recommended for security + sliding sessions)
	newRefresh, newRefreshExp, err := s.token.GenerateRefreshToken(dbUserID)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if err := s.sessionStore.UpdateTokens(ctx, sess.ID, newAccess, newAccessExp, newRefresh, newRefreshExp); err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	sess.AccessToken = newAccess
	sess.AccessTokenExp = newAccessExp
	sess.RefreshToken = newRefresh
	sess.RefreshTokenExp = newRefreshExp
	sess.LastUsedAt = &now
	return sess, nil
}

func (s *SessionUsecase) Logout(ctx context.Context, sessionID uint64) error {
	if err := s.sessionStore.DeleteByID(ctx, sessionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (s *SessionUsecase) LogoutAll(ctx context.Context, userID uint64) error {
	if err := s.sessionStore.DeleteByUserID(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (s *SessionUsecase) EnforceMaxDevices(ctx context.Context, userID uint64, max int) error {
	if max <= 0 {
		max = 5
	}

	sessions, err := s.sessionStore.GetAllValidSessionsByUserId(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	for len(sessions) > max {
		if err := s.sessionStore.DeleteOldestValidSession(ctx, userID); err != nil {
			return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
		sessions = sessions[:len(sessions)-1]
	}
	return nil
}

func (s *SessionUsecase) GetSessionsByUserID(ctx context.Context, userID uint64) ([]domain.UserSession, error) {
	return s.sessionStore.GetAllValidSessionsByUserId(ctx, userID)
}

func (s *SessionUsecase) RevokeSessionByID(ctx context.Context, sessionID, userID uint64) error {
	return s.sessionStore.RevokeByID(ctx, sessionID, userID)
}
