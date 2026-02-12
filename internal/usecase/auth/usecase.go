package authUsecase

import (
	"context"
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
		err error
	)

	isEmail := strings.Contains(req.LoginInput, "@")

	if isEmail {
		user, err = a.authStore.GetByEmail(ctx, req.LoginInput)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to get user by email")
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
	}

	if !isEmail {
		user, err = a.authStore.GetByPhone(ctx, req.LoginInput)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to get user by phone")
			return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
		}
	}

	if user == nil {
		return nil, apperr.New(apperr.CodeNotFound, http.StatusNotFound, "user not found")
	}

	if !user.Verified {
		return nil, apperr.New(apperr.CodeConflict, http.StatusConflict, "user is not verified")
	}

	if err := a.hasher.CheckPasswordHash(req.Password, user.Password); err != nil {
		return nil, apperr.New(apperr.CodeUnauthorized, http.StatusUnauthorized, "invalid password")
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

	session := &domain.UserSession{
		UserID:          user.ID,
		AccessToken:     accessToken,
		AccessTokenExp:  accessTokenExp,
		RefreshToken:    refreshToken,
		RefreshTokenExp: refreshTokenExp,
		IPAddress:       meta.IP,
		Device:          meta.Device,
		UserAgent:       meta.UserAgent,
	}

	// create session
	if err := a.session.Create(ctx, session); err != nil {
		a.logger.Error().Err(err).Msg("failed to create session")
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	resp := &LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		AccessTokenExp: accessTokenExp.String(),
		RefreshTokenExp: refreshTokenExp.String(),
		Device: meta.Device,
		UserEmail: user.Email,
		IpAddress: meta.IP,
	}

	return resp, nil
}

func (a *AuthUsecase) Refresh(ctx context.Context, req RefreshTokenRequest, meta SessionMeta) (*RefreshTokenResponse, error) {

	return nil, nil
}

// helpers
func generateRandomCode() int {
	min := 100000
	max := 999999

	return rand.IntN(max-min+1) + min
}
