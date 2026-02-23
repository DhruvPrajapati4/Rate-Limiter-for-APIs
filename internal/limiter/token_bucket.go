package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	redisscripts "github.com/DhruvPrajapati4/rate-limiter/internal/redis"
	"github.com/redis/go-redis/v9"
)

// TokenBucket implements the token bucket rate limiting algorithm backed by Redis.
type TokenBucket struct {
	rdb        *redis.Client
	burstSize  int
	refillRate float64 // tokens per second
	window     time.Duration
	limit      int
}

// NewTokenBucket creates a new token bucket limiter.
func NewTokenBucket(cfg config.TierDetail, rdb *redis.Client) *TokenBucket {
	burstSize := cfg.BurstSize
	if burstSize == 0 {
		burstSize = cfg.Limit
	}
	refillRate := float64(cfg.RefillRate)
	if cfg.RefillRate == 0 {
		refillRate = float64(cfg.Limit) / cfg.Window.Seconds()
	}

	return &TokenBucket{
		rdb:        rdb,
		burstSize:  burstSize,
		refillRate: refillRate,
		window:     cfg.Window,
		limit:      cfg.Limit,
	}
}

// Allow checks if a request is allowed under the token bucket algorithm.
func (tb *TokenBucket) Allow(ctx context.Context, key string) (Result, error) {
	now := float64(time.Now().UnixNano()) / 1e9
	ttl := int(tb.window.Seconds()) * 2 // TTL is 2x window for safety

	vals, err := redisscripts.TokenBucketScript.Run(ctx, tb.rdb,
		[]string{fmt.Sprintf("tb:%s", key)},
		tb.burstSize,
		tb.refillRate,
		now,
		ttl,
	).Int64Slice()
	if err != nil {
		return Result{}, fmt.Errorf("token bucket script: %w", err)
	}

	allowed := vals[0] == 1
	remaining := int(vals[1])

	result := Result{
		Allowed:   allowed,
		Remaining: remaining,
		Limit:     tb.limit,
		ResetAt:   time.Now().Add(tb.window),
	}

	if !allowed {
		// Time until at least 1 token is available
		result.RetryAfter = time.Duration(float64(time.Second) / tb.refillRate)
	}

	return result, nil
}
