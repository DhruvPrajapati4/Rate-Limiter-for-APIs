package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Server    ServerConfig    `yaml:"server"`
		Redis     RedisConfig     `yaml:"redis"`
		RateLimit RateLimitConfig `yaml:"rate_limit"`
	}

	ServerConfig struct {
		Port int `yaml:"port"`
	}

	RedisConfig struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}

	RateLimitConfig struct {
		Algorithm string     `yaml:"algorithm"`
		Tiers     TierConfig `yaml:"tiers"`
	}

	TierConfig struct {
		APISpecific TierDetail `yaml:"api_specific"`
		GlobalUser  TierDetail `yaml:"global_user"`
		GlobalAPI   TierDetail `yaml:"global_api"`
	}

	TierDetail struct {
		Limit      int           `yaml:"limit"`
		Window     time.Duration `yaml:"window"`
		BurstSize  int           `yaml:"burst_size,omitempty"`
		RefillRate int           `yaml:"refill_rate,omitempty"`
	}
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Allow environment variable overrides
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.Redis.Addr = addr
	}

	return &cfg, nil
}
