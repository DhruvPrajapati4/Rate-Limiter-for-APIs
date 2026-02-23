package limiter

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	redisscripts "github.com/DhruvPrajapati4/rate-limiter/internal/redis"
	"github.com/redis/go-redis/v9"
)

// SlidingWindow implements the sliding window counter rate limiting algorithm backed by Redis.
type SlidingWindow struct {
	rdb    *redis.Client
	limit  int
	window time.Duration
}

// NewSlidingWindow creates a new sliding window counter limiter.
func NewSlidingWindow(cfg config.TierDetail, rdb *redis.Client) *SlidingWindow {
	return &SlidingWindow{
		rdb:    rdb,
		limit:  cfg.Limit,
		window: cfg.Window,
	}
}

// Allow checks if a request is allowed under the sliding window counter algorithm.
func (sw *SlidingWindow) Allow(ctx context.Context, key string) (Result, error) {
	now := float64(time.Now().UnixNano()) / 1e9
	windowSecs := sw.window.Seconds()
	ttl := int(windowSecs) * 2

	// Calculate current and previous window keys
	currWindow := int64(math.Floor(now / windowSecs))
	currKey := fmt.Sprintf("sw:%s:%d", key, currWindow)
	prevKey := fmt.Sprintf("sw:%s:%d", key, currWindow-1)

	vals, err := redisscripts.SlidingWindowScript.Run(ctx, sw.rdb,
		[]string{currKey, prevKey},
		sw.limit,
		windowSecs,
		now,
		ttl,
	).Int64Slice()
	if err != nil {
		return Result{}, fmt.Errorf("sliding window script: %w", err)
	}

	allowed := vals[0] == 1
	remaining := int(vals[1])

	result := Result{
		Allowed:   allowed,
		Remaining: remaining,
		Limit:     sw.limit,
		ResetAt:   time.Unix(int64((currWindow+1)*int64(windowSecs)), 0),
	}

	if !allowed {
		resetTime := float64((currWindow + 1)) * windowSecs
		result.RetryAfter = time.Duration((resetTime - now) * float64(time.Second))
	}

	return result, nil
}
