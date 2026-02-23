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
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		DB           int    `yaml:"db"`
		Password     string `yaml:"password"`
		Username     string `yaml:"username"`
		PoolSize     int    `yaml:"poolSize"`
		MinIdleConns int    `yaml:"minIdleConns"`
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

	password := os.Getenv("REDIS_PASSWORD")
	if password != "" {
		cfg.Redis.Password = password
	}

	return &cfg, nil
}
