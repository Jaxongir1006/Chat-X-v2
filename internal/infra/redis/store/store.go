package redisStore

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPRedisStore struct {
	rdb *redis.Client
}

func NewOTPRedisStore(rdb *redis.Client) *OTPRedisStore {
	return &OTPRedisStore{rdb: rdb}
}

func (s *OTPRedisStore) key(email string) string {
	return "otp:email:" + email
}

func (s *OTPRedisStore) SaveEmailCode(ctx context.Context, email, codeHash string, ttl time.Duration) error {
	return s.rdb.Set(ctx, s.key(email), codeHash, ttl).Err()
}

func (s *OTPRedisStore) GetEmailCodeHash(ctx context.Context, email string) (string, error) {
	val, err := s.rdb.Get(ctx, s.key(email)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *OTPRedisStore) DeleteEmailCode(ctx context.Context, email string) error {
	return s.rdb.Del(ctx, s.key(email)).Err()
}
