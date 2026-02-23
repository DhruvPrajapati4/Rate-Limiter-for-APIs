package throttle

import (
	"context"
	"fmt"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	"github.com/DhruvPrajapati4/rate-limiter/internal/limiter"
	"github.com/redis/go-redis/v9"
)

type (
	// MultiTier composes three rate limiters (API-specific, global user, global API)
	// and checks all tiers concurrently using goroutines and channels.
	MultiTier struct {
		apiLimiter    limiter.Limiter
		userLimiter   limiter.Limiter
		globalLimiter limiter.Limiter
	}

	// tierResult holds the outcome of a single tier check.
	tierResult struct {
		tier   string
		result limiter.Result
		err    error
	}
)

// NewMultiTier creates a MultiTier throttler with the given config and Redis client.
func NewMultiTier(cfg config.RateLimitConfig, rdb *redis.Client) *MultiTier {
	return &MultiTier{
		apiLimiter:    limiter.New(cfg.Algorithm, cfg.Tiers.APISpecific, rdb),
		userLimiter:   limiter.New(cfg.Algorithm, cfg.Tiers.GlobalUser, rdb),
		globalLimiter: limiter.New(cfg.Algorithm, cfg.Tiers.GlobalAPI, rdb),
	}
}

// Allow checks all three tiers concurrently and returns the most restrictive result.
// If any tier denies the request, the overall result is denied.
func (m *MultiTier) Allow(ctx context.Context, userID, apiName string) (limiter.Result, error) {
	ch := make(chan tierResult, 3)

	// Launch all three tier checks concurrently
	go func() {
		r, err := m.apiLimiter.Allow(ctx, fmt.Sprintf("api:%s:user:%s", apiName, userID))
		ch <- tierResult{"api_specific", r, err}
	}()
	go func() {
		r, err := m.userLimiter.Allow(ctx, fmt.Sprintf("user:%s", userID))
		ch <- tierResult{"global_user", r, err}
	}()
	go func() {
		r, err := m.globalLimiter.Allow(ctx, fmt.Sprintf("global_api:%s", apiName))
		ch <- tierResult{"global_api", r, err}
	}()

	// Collect results - return the most restrictive
	var final limiter.Result
	final.Allowed = true

	for range 3 {
		tr := <-ch
		if tr.err != nil {
			return limiter.Result{}, fmt.Errorf("tier %s: %w", tr.tier, tr.err)
		}
		if !tr.result.Allowed {
			final = tr.result
			final.Allowed = false
		} else if final.Allowed && tr.result.Remaining < final.Remaining {
			// Among allowed results, track the one with fewest remaining
			final = tr.result
		}
	}

	return final, nil
}
