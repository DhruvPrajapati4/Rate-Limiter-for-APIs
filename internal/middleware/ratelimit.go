package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DhruvPrajapati4/rate-limiter/internal/throttle"
	"github.com/gin-gonic/gin"
)

// RateLimit returns a Gin middleware that enforces multi-tier rate limiting.
// It extracts the user ID and API name from the request, runs the throttler,
// and sets standard rate limit response headers.
func RateLimit(throttler *throttle.MultiTier) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := extractUserID(c)
		apiName := c.FullPath()
		if apiName == "" {
			apiName = c.Request.URL.Path
		}

		result, err := throttler.Allow(c.Request.Context(), userID, apiName)
		if err != nil {
			// Fail-open: allow the request through but log the error.
			// This prevents a Redis outage from causing a full service outage.
			log.Printf("rate limiter error (fail-open): %v", err)
			c.Next()
			return
		}

		// Set rate limit headers on every response
		c.Header("X-RateLimit-Limit", strconv.Itoa(result.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

		if !result.Allowed {
			retryAfter := int(result.RetryAfter.Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": result.RetryAfter.Seconds(),
			})
			return
		}

		c.Next()
	}
}

// extractUserID extracts the user identifier from the request.
// Priority: path param "userId" > X-User-ID header > query param "userId" > client IP.
func extractUserID(c *gin.Context) string {
	if id := c.Param("userId"); id != "" {
		return id
	}
	if id := c.GetHeader("X-User-ID"); id != "" {
		return id
	}
	if id := c.Query("userId"); id != "" {
		return id
	}
	return c.ClientIP()
}
