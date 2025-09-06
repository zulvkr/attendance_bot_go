package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all application configuration
type Config struct {
	BotToken      string
	TOTPSecret    string
	AdminPassword string
	Environment   string
	DatabasePath  string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		TOTPSecret:    os.Getenv("TOTP_SECRET"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		Environment:   getEnvWithDefault("NODE_ENV", "development"),
		DatabasePath:  getEnvWithDefault("DATABASE_PATH", "data/attendance.db"),
	}

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// validate ensures all required configuration is present
func (c *Config) validate() error {
	var missing []string

	if c.BotToken == "" {
		missing = append(missing, "BOT_TOKEN")
	}
	if len(c.BotToken) < 10 {
		missing = append(missing, "BOT_TOKEN (must be at least 10 characters)")
	}

	if c.TOTPSecret == "" {
		missing = append(missing, "TOTP_SECRET")
	}
	if len(c.TOTPSecret) < 16 {
		missing = append(missing, "TOTP_SECRET (must be at least 16 characters)")
	}

	if c.AdminPassword == "" {
		missing = append(missing, "ADMIN_PASSWORD")
	}
	if len(c.AdminPassword) < 8 {
		missing = append(missing, "ADMIN_PASSWORD (must be at least 8 characters)")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing or invalid environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnvWithDefault returns the environment variable value or a default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
