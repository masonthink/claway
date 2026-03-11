package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DatabaseURL          string
	RedisURL             string
	Port                 string
	JWTSecret            string
	OpenClawClientID     string
	OpenClawClientSecret string
	OpenClawBaseURL      string
	UpstreamLLMBaseURL   string
	UpstreamLLMAPIKey    string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		RedisURL:             os.Getenv("REDIS_URL"),
		Port:                 os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		OpenClawClientID:     os.Getenv("OPENCLAW_CLIENT_ID"),
		OpenClawClientSecret: os.Getenv("OPENCLAW_CLIENT_SECRET"),
		OpenClawBaseURL:      os.Getenv("OPENCLAW_BASE_URL"),
		UpstreamLLMBaseURL:   os.Getenv("UPSTREAM_LLM_BASE_URL"),
		UpstreamLLMAPIKey:    os.Getenv("UPSTREAM_LLM_API_KEY"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// Redis is optional for MVP

	if cfg.OpenClawBaseURL == "" {
		cfg.OpenClawBaseURL = "https://api.openclaw.ai"
	}

	if cfg.UpstreamLLMBaseURL == "" {
		cfg.UpstreamLLMBaseURL = "https://api.openai.com"
	}

	return cfg, nil
}
