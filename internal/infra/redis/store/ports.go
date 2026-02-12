package redisStore

import (
	"context"
	"time"
)

type OTPStore interface {
	SaveEmailCode(ctx context.Context, email string, codeHash string, ttl time.Duration) error
	GetEmailCodeHash(ctx context.Context, email string) (string, error)
	DeleteEmailCode(ctx context.Context, email string) error
}
