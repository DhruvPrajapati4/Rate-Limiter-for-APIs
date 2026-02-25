package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	goredis "github.com/redis/go-redis/v9"
)

// NewClient creates a new Redis client and verifies the connection.
func NewClient(cfg config.RedisConfig) (*goredis.Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return rdb, nil
}
