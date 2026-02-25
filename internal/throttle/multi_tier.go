package throttle

import (
	"context"
	"sync"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	"github.com/DhruvPrajapati4/rate-limiter/internal/limiter"
	"github.com/redis/go-redis/v9"
)

type (
	// MultiTier composes three rate limiters (API-specific, global user, global API)
	// and checks all tiers concurrently using goroutines.
	// The most restrictive result is returned. If any tier denies, the request is denied.
	MultiTier struct {
		apiLimiter    limiter.Limiter
		userLimiter   limiter.Limiter
		globalLimiter limiter.Limiter
	}

	tierResult struct {
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

// Allow checks all three tiers concurrently using goroutines and returns the
// most restrictive result. If any tier denies, the request is denied.
//
// Note: Because all tiers are evaluated in parallel, tokens are consumed from
// all tiers even if one denies (phantom consumption). This is an acceptable
// trade-off for lower latency — tokens refill naturally, and denial only
// occurs when the user is already at or near their limit.
func (m *MultiTier) Allow(ctx context.Context, userID, apiName string) (limiter.Result, error) {
	type tierCheck struct {
		limiter limiter.Limiter
		key     string
	}

	checks := []tierCheck{
		{m.apiLimiter, "api:" + apiName + ":user:" + userID},
		{m.userLimiter, "user:" + userID},
		{m.globalLimiter, "global_api:" + apiName},
	}

	results := make([]tierResult, len(checks))
	var wg sync.WaitGroup
	wg.Add(len(checks))

	for i, tc := range checks {
		go func(idx int, l limiter.Limiter, key string) {
			defer wg.Done()
			r, err := l.Allow(ctx, key)
			results[idx] = tierResult{result: r, err: err}
		}(i, tc.limiter, tc.key)
	}

	wg.Wait()

	// Collect results: return first error, find most restrictive
	var mostRestrictive limiter.Result
	first := true

	for _, tr := range results {
		if tr.err != nil {
			return limiter.Result{}, tr.err
		}

		if first {
			mostRestrictive = tr.result
			first = false
		} else if !tr.result.Allowed {
			mostRestrictive = tr.result
		} else if mostRestrictive.Allowed && tr.result.Remaining < mostRestrictive.Remaining {
			mostRestrictive = tr.result
		}
	}

	return mostRestrictive, nil
}
