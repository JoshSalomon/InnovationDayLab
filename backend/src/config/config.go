package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port        string
	DatabasePath string
	SessionKey  string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "./data/tasks.db"),
		SessionKey:  getEnv("SESSION_KEY", "change-me-in-production-secret-key-min-32-chars"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
