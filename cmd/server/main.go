package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/DhruvPrajapati4/rate-limiter/internal/config"
	"github.com/DhruvPrajapati4/rate-limiter/internal/middleware"
	redisclient "github.com/DhruvPrajapati4/rate-limiter/internal/redis"
	"github.com/DhruvPrajapati4/rate-limiter/internal/throttle"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rdb, err := redisclient.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("warning: failed to close redis: %v", err)
		}
	}()

	log.Printf("config loaded: algorithm=%s", cfg.RateLimit.Algorithm)
	log.Printf("  api_specific: limit=%d, window=%s, burst=%d, refill=%d",
		cfg.RateLimit.Tiers.APISpecific.Limit, cfg.RateLimit.Tiers.APISpecific.Window,
		cfg.RateLimit.Tiers.APISpecific.BurstSize, cfg.RateLimit.Tiers.APISpecific.RefillRate)
	log.Printf("  global_user:  limit=%d, window=%s, burst=%d, refill=%d",
		cfg.RateLimit.Tiers.GlobalUser.Limit, cfg.RateLimit.Tiers.GlobalUser.Window,
		cfg.RateLimit.Tiers.GlobalUser.BurstSize, cfg.RateLimit.Tiers.GlobalUser.RefillRate)
	log.Printf("  global_api:   limit=%d, window=%s, burst=%d, refill=%d",
		cfg.RateLimit.Tiers.GlobalAPI.Limit, cfg.RateLimit.Tiers.GlobalAPI.Window,
		cfg.RateLimit.Tiers.GlobalAPI.BurstSize, cfg.RateLimit.Tiers.GlobalAPI.RefillRate)

	mt := throttle.NewMultiTier(cfg.RateLimit, rdb)

	r := gin.Default()
	r.Use(middleware.RateLimit(mt))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Sample APIs for testing rate limiter
	api := r.Group("/api/:userId")
	{
		api.GET("/alpha", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "API-ALPHA response", "user": c.Param("userId")})
		})
		api.POST("/beta", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "API-BETA response", "user": c.Param("userId")})
		})
		api.PUT("/gamma", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "API-GAMMA response", "user": c.Param("userId")})
		})
		api.DELETE("/delta", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "API-DELTA response", "user": c.Param("userId")})
		})
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("starting rate limiter server on %s (algorithm: %s)", addr, cfg.RateLimit.Algorithm)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
