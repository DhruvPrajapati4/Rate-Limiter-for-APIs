package limiter

import (
	"context"
	"time"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	"github.com/redis/go-redis/v9"
)

// Result holds the outcome of a rate limit check.
type (
	Result struct {
		Allowed    bool
		Remaining  int
		Limit      int
		RetryAfter time.Duration
		ResetAt    time.Time
	}

	// Limiter defines the interface for rate limiting algorithms.
	Limiter interface {
		Allow(ctx context.Context, key string) (Result, error)
	}
)

// New creates a Limiter based on the configured algorithm.
func New(algorithm string, cfg config.TierDetail, rdb *redis.Client) Limiter {
	switch algorithm {
	case "sliding_window": // TODO: add constants
		return NewSlidingWindow(cfg, rdb)
	default:
		return NewTokenBucket(cfg, rdb)
	}
}
