package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalPort := os.Getenv("PORT")
	originalDBPath := os.Getenv("DATABASE_PATH")
	originalSessionKey := os.Getenv("SESSION_KEY")

	// Clean up after test
	defer func() {
		if originalPort != "" {
			os.Setenv("PORT", originalPort)
		} else {
			os.Unsetenv("PORT")
		}
		if originalDBPath != "" {
			os.Setenv("DATABASE_PATH", originalDBPath)
		} else {
			os.Unsetenv("DATABASE_PATH")
		}
		if originalSessionKey != "" {
			os.Setenv("SESSION_KEY", originalSessionKey)
		} else {
			os.Unsetenv("SESSION_KEY")
		}
	}()

	t.Run("default values", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_PATH")
		os.Unsetenv("SESSION_KEY")

		cfg := Load()

		if cfg.Port != "8080" {
			t.Errorf("Load() Port = %v, want 8080", cfg.Port)
		}
		if cfg.DatabasePath != "./data/tasks.db" {
			t.Errorf("Load() DatabasePath = %v, want ./data/tasks.db", cfg.DatabasePath)
		}
		if cfg.SessionKey != "change-me-in-production-secret-key-min-32-chars" {
			t.Errorf("Load() SessionKey = %v, want change-me-in-production-secret-key-min-32-chars", cfg.SessionKey)
		}
	})

	t.Run("environment variables", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		os.Setenv("DATABASE_PATH", "/custom/path/db.db")
		os.Setenv("SESSION_KEY", "custom-secret-key")

		cfg := Load()

		if cfg.Port != "3000" {
			t.Errorf("Load() Port = %v, want 3000", cfg.Port)
		}
		if cfg.DatabasePath != "/custom/path/db.db" {
			t.Errorf("Load() DatabasePath = %v, want /custom/path/db.db", cfg.DatabasePath)
		}
		if cfg.SessionKey != "custom-secret-key" {
			t.Errorf("Load() SessionKey = %v, want custom-secret-key", cfg.SessionKey)
		}
	})

	t.Run("partial environment variables", func(t *testing.T) {
		os.Setenv("PORT", "9000")
		os.Unsetenv("DATABASE_PATH")
		os.Unsetenv("SESSION_KEY")

		cfg := Load()

		if cfg.Port != "9000" {
			t.Errorf("Load() Port = %v, want 9000", cfg.Port)
		}
		if cfg.DatabasePath != "./data/tasks.db" {
			t.Errorf("Load() DatabasePath = %v, want ./data/tasks.db (default)", cfg.DatabasePath)
		}
	})
}

func TestGetEnv(t *testing.T) {
	originalValue := os.Getenv("TEST_VAR")

	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_VAR")
		}
	}()

	t.Run("returns environment value when set", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test-value")
		got := getEnv("TEST_VAR", "default")
		if got != "test-value" {
			t.Errorf("getEnv() = %v, want test-value", got)
		}
	})

	t.Run("returns default when not set", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		got := getEnv("TEST_VAR", "default-value")
		if got != "default-value" {
			t.Errorf("getEnv() = %v, want default-value", got)
		}
	})

	t.Run("returns default when empty", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		got := getEnv("TEST_VAR", "default-value")
		if got != "default-value" {
			t.Errorf("getEnv() = %v, want default-value (empty string should use default)", got)
		}
	})
}
