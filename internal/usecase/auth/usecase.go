package authUsecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/redis/go-redis/v9"
)

func (a *AuthUsecase) Register(ctx context.Context, req RegisterRequest) error {
	if req.ConfirmPass != req.Password {
		return apperr.New(apperr.CodeConflict, http.StatusConflict, "passwords do not match")
	}

	// password validation
	if len(req.Password) < 8 {
		return apperr.New(apperr.CodeConflict, http.StatusConflict, "password must be at least 8 characters long")
	}

	hashed, err := a.hasher.Hash(req.Password)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// create user
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashed,
		Role:     "user",
		Verified: false,
	}
	if err := a.authStore.InsertUser(ctx, user); err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	go func() {
		// save the email code to redis with a 5 minute ttl
		bgCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		code := generateRandomCode()
		hashedCode := a.codeHasher.Hash(fmt.Sprintf("%d", code))
		err = a.redis.SaveEmailCode(bgCtx, req.Email, hashedCode, time.Minute*5)
		if err != nil {
			a.logger.Error().Err(err).Str("email", req.Email).Int("code", code).Msg("failed to save email code to redis")
			return
		}
		a.logger.Debug().Str("email", req.Email).Int("code", code).Msg("email code saved to redis")
		// send email to user with the email code.
	}()

	return nil
}

func (a *AuthUsecase) VerifyUser(ctx context.Context, email string, code int, meta SessionMeta) (*VerifyUserResponse, error) {
	codeHash, err := a.redis.GetEmailCodeHash(ctx, email)
	if err != nil {
		if err == redis.Nil {
			return nil, apperr.New(apperr.CodeConflict, http.StatusConflict, "email code is invalid")
		}
		a.logger.Error().Err(err).Msg("failed to get hashed code from redis")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if ok := a.codeHasher.Compare(fmt.Sprintf("%d", code), codeHash); !ok {
		a.logger.Error().Str("email", email).Int("code", code).Msg("email code is invalid")
		return nil, apperr.New(apperr.CodeConflict, http.StatusConflict, "email code is invalid")
	}

	err = a.redis.DeleteEmailCode(ctx, email)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to delete email code from redis")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// get user by email
	user, err := a.authStore.GetByEmail(ctx, email)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to get user by email")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// create tokens
	accessToken, accessTokenExp, err := a.token.GenerateAccessToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate access token")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	refreshToken, refreshTokenExp, err := a.token.GenerateRefreshToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate refresh token")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// create session for the user
	session := &domain.UserSession{
		UserID:          user.ID,
		AccessToken:     accessToken,
		AccessTokenExp:  accessTokenExp,
		RefreshToken:    refreshToken,
		RefreshTokenExp: refreshTokenExp,
	}

	if err := a.session.Create(ctx, session); err != nil {
		a.logger.Error().Err(err).Msg("failed to create session")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	err = a.authStore.VerifyUser(ctx, email)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to verify user")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	resp := &VerifyUserResponse{
		AccessToken:     accessToken,
		AccessTokenExp:  accessTokenExp.String(),
		RefreshToken:    refreshToken,
		RefreshTokenExp: refreshTokenExp.String(),
		UserEmail:       email,
		IpAddress:       meta.IP,
		Device:          meta.Device,
	}

	return resp, nil
}

func (a *AuthUsecase) Login(ctx context.Context, req LoginRequest, meta SessionMeta) (*LoginResponse, error) {
	var (
		user *domain.User
		err  error
	)

	isEmail := strings.Contains(req.LoginInput, "@")

	if isEmail {
		user, err = a.authStore.GetByEmail(ctx, req.LoginInput)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "invalid credentials")
			}
			a.logger.Error().Err(err).Msg("failed to get user by email")
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
	} else {
		user, err = a.authStore.GetByPhone(ctx, req.LoginInput)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "invalid credentials")
			}
			a.logger.Error().Err(err).Msg("failed to get user by phone")
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
	}

	if user == nil {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "invalid credentials")
	}

	if !user.Verified {
		return nil, apperr.New(apperr.CodeConflict, http.StatusConflict, "user is not verified")
	}

	if err := a.hasher.CheckPasswordHash(req.Password, user.Password); err != nil {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "invalid credentials")
	}

	accessToken, accessTokenExp, err := a.token.GenerateAccessToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate access token")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	refreshToken, refreshTokenExp, err := a.token.GenerateRefreshToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate refresh token")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	sessions, err := a.session.GetAllValidSessionsByUserId(ctx, user.ID)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to get user sessions")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	// If already at/over limit, delete oldest until there's room (at least once here).
	// If you want to be strict, loop while len(sessions) >= 5.
	if len(sessions) >= 5 {
		if err := a.session.DeleteOldestValidSession(ctx, user.ID); err != nil {
			a.logger.Error().Err(err).Msg("failed to delete oldest session")
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
		// Optional: you can re-fetch sessions here if your Delete affects which one should be updated.
	}

	updated := false
	for _, s := range sessions {
		// NOTE: Device string matching can be weak. Prefer meta.DeviceID if you have it.
		if s.Device == meta.Device {
			updated = true

			// Update tokens for this existing session
			if err := a.session.UpdateTokens(ctx, s.ID, accessToken, accessTokenExp, refreshToken, refreshTokenExp); err != nil {
				a.logger.Error().Err(err).Msg("failed to update session tokens")
				return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
			}

			// Keep meta fresh too (recommended). If you don't have UpdateMeta, add it.
			// If you already update meta inside UpdateTokens, remove this block.
			if err := a.session.UpdateMeta(ctx, s.ID, meta.Device, meta.IP, meta.UserAgent, time.Now()); err != nil {
				a.logger.Error().Err(err).Msg("failed to update session meta")
				return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
			}

			break
		}
	}

	now := time.Now()

	if !updated {
		session := &domain.UserSession{
			UserID:          user.ID,
			AccessToken:     accessToken,
			AccessTokenExp:  accessTokenExp,
			RefreshToken:    refreshToken,
			RefreshTokenExp: refreshTokenExp,
			IPAddress:       meta.IP,
			Device:          meta.Device,
			UserAgent:       meta.UserAgent,
			LastUsedAt:      &now,
		}

		if err := a.session.Create(ctx, session); err != nil {
			a.logger.Error().Err(err).Msg("failed to create session")
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
	}

	return &LoginResponse{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessTokenExp:  accessTokenExp.Format(time.RFC3339),
		RefreshTokenExp: refreshTokenExp.Format(time.RFC3339),
		Device:          meta.Device,
		UserEmail:       user.Email,
		IpAddress:       meta.IP,
	}, nil
}

func (a *AuthUsecase) Refresh(ctx context.Context, req RefreshTokenRequest, meta SessionMeta) (*RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", errors.New("missing refresh token"))
	}

	claims, err := a.token.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", err)
	}

	sess, err := a.session.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", err)
		}
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if sess.RevokedAt != nil {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", errors.New("session revoked"))
	}

	now := time.Now()
	if !sess.RefreshTokenExp.After(now) {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", errors.New("refresh session expired"))
	}

	// CRITICAL: bind claims -> session user.
	// Adjust accessor: Subject / UserID / GetSubject etc.
	claimSub := ""
	// Example possibilities (uncomment the one matching your claims type):
	claimSub = claims.Subject
	// claimSub = claims.GetSubject()
	if claimSub == "" {
		// If you can't access subject, at least log it and fail safe.
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", errors.New("invalid refresh claims subject"))
	}

	if claimSub != fmt.Sprintf("%d", sess.UserID) {
		return nil, apperr.Wrap(apperr.CodeUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED", errors.New("token user mismatch"))
	}

	accessToken, accessExp, err := a.token.GenerateAccessToken(fmt.Sprintf("%d", sess.UserID))
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate access token")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	newRefresh, newRefreshExp, err := a.token.GenerateRefreshToken(fmt.Sprintf("%d", sess.UserID))
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate refresh token")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if err := a.session.RotateRefresh(ctx, sess.ID, newRefresh, newRefreshExp); err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	if err := a.session.UpdateMeta(ctx, sess.ID, meta.Device, meta.IP, meta.UserAgent, now); err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	return &RefreshTokenResponse{
		AccessToken:     accessToken,
		AccessTokenExp:  accessExp.Format(time.RFC3339),
		RefreshToken:    newRefresh,
		RefreshTokenExp: newRefreshExp.Format(time.RFC3339),
		Device:          meta.Device,
		IpAddress:       meta.IP,
	}, nil
}

func (a *AuthUsecase) LogoutFromCurrent(ctx context.Context, userID uint64, sessionID uint64) error {
	err := a.session.RevokeByID(ctx, sessionID, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (a *AuthUsecase) LogoutAll(ctx context.Context, userID uint64) error {
	err := a.session.RevokeAllByUserID(ctx, userID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

func (a *AuthUsecase) LogOutAllExceptCurrent(ctx context.Context, userID uint64, sessionID uint64) error {
	err := a.session.RevokeAllExceptCurrent(ctx, userID, sessionID)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	return nil
}

// helpers
func generateRandomCode() int {
	min := 100000
	max := 999999

	return rand.IntN(max-min+1) + min
}
