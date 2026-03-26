package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DatabaseURL        string
	RedisURL           string
	Port               string
	JWTSecret          string
	FrontendURL        string
	UpstreamLLMBaseURL string
	UpstreamLLMAPIKey  string

	// X (Twitter) OAuth 2.0
	XClientID     string
	XClientSecret string
	XRedirectURI  string // backend callback URL

	// GitHub OAuth 2.0
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURI  string

	// Google OAuth 2.0
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string

	// Legacy OpenClaw OAuth (kept for backward compatibility)
	OpenClawClientID     string
	OpenClawClientSecret string
	OpenClawBaseURL      string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		RedisURL:             os.Getenv("REDIS_URL"),
		Port:                 os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		FrontendURL:          os.Getenv("FRONTEND_URL"),
		UpstreamLLMBaseURL:   os.Getenv("UPSTREAM_LLM_BASE_URL"),
		UpstreamLLMAPIKey:    os.Getenv("UPSTREAM_LLM_API_KEY"),
		XClientID:            os.Getenv("X_CLIENT_ID"),
		XClientSecret:        os.Getenv("X_CLIENT_SECRET"),
		XRedirectURI:         os.Getenv("X_REDIRECT_URI"),
		GitHubClientID:       os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubRedirectURI:    os.Getenv("GITHUB_REDIRECT_URI"),
		GoogleClientID:       os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:    os.Getenv("GOOGLE_REDIRECT_URI"),
		OpenClawClientID:     os.Getenv("OPENCLAW_CLIENT_ID"),
		OpenClawClientSecret: os.Getenv("OPENCLAW_CLIENT_SECRET"),
		OpenClawBaseURL:      os.Getenv("OPENCLAW_BASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters for security")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.FrontendURL == "" {
		cfg.FrontendURL = "https://claway.cc"
	}

	if cfg.XRedirectURI == "" {
		cfg.XRedirectURI = "https://api.claway.cc/api/v1/auth/x/callback"
	}

	if cfg.GitHubRedirectURI == "" {
		cfg.GitHubRedirectURI = "https://api.claway.cc/api/v1/auth/github/callback"
	}

	if cfg.GoogleRedirectURI == "" {
		cfg.GoogleRedirectURI = "https://api.claway.cc/api/v1/auth/google/callback"
	}

	if cfg.OpenClawBaseURL == "" {
		cfg.OpenClawBaseURL = "https://api.openclaw.ai"
	}

	if cfg.UpstreamLLMBaseURL == "" {
		cfg.UpstreamLLMBaseURL = "https://api.openai.com"
	}

	return cfg, nil
}
