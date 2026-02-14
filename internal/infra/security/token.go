package security

import (
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenStore interface {
	GenerateAccessToken(userID string) (string, time.Time, error)
	GenerateRefreshToken(userID string) (string, time.Time, error)
	VerifyAccessToken(tokenStr string) (*Claims, error)
	VerifyRefreshToken(tokenStr string) (*Claims, error)
	GetUserIDFromAccess(tokenStr string) (string, error)
}

type Token struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func NewToken(cfg config.TokenConfig) *Token {
	return &Token{
		AccessSecret:  cfg.AccessSecret,
		RefreshSecret: cfg.RefreshSecret,
		AccessTTL:     cfg.AccessTTL,
		RefreshTTL:    cfg.RefreshTTL,
	}
}

func (t *Token) GenerateAccessToken(userID string) (string, time.Time, error) {
	exp := time.Now().Add(t.AccessTTL)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(t.AccessSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, exp, nil
}

func (t *Token) GenerateRefreshToken(userID string) (string, time.Time, error) {
	exp := time.Now().Add(t.RefreshTTL)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(t.RefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, exp, nil
}

var ErrInvalidToken = apperr.New(apperr.CodeUnauthorized, 401, "UNAUTHORIZED")

func (t *Token) VerifyAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.AccessSecret), nil
	})

	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (t *Token) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.RefreshSecret), nil
	})

	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (t *Token) GetUserIDFromAccess(tokenStr string) (string, error) {
	claims, err := t.VerifyAccessToken(tokenStr)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
