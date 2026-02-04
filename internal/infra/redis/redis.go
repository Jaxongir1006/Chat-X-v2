package redisInfra

import (
	"context"
	"fmt"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	Config config.RedisConfig
	Client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) *redisClient {
	return &redisClient{
		Config: cfg,
	}
}

func (r *redisClient) InitRedis() error {
	if r.Config.Host == "" {
		return fmt.Errorf("redis host is empty")
	}

	redisUrl := fmt.Sprintf("%s:%d", r.Config.Host, r.Config.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: r.Config.Pass,
		DB:       0,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

func (r *redisClient) Close() error {
	if r == nil || r.Client == nil {
		return nil
	}
	return r.Client.Close()
}
