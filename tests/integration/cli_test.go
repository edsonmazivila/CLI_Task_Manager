package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/edson-mazvila/task-manager/internal/config"
)

// TestCLIConfiguration tests configuration loading
func TestCLIConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "valid_sqlite_config",
			envVars: map[string]string{
				"DB_TYPE": "sqlite",
				"DB_PATH": "/tmp/test.db",
			},
			expectError: false,
		},
		{
			name: "invalid_db_type",
			envVars: map[string]string{
				"DB_TYPE": "invalid",
			},
			expectError: true,
		},
		{
			name: "missing_required_postgres_fields",
			envVars: map[string]string{
				"DB_TYPE": "postgres",
				"DB_HOST": "localhost",
			},
			expectError: true,
		},
		{
			name:    "default_values",
			envVars: map[string]string{
				// No env vars, should use defaults
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			cfg, err := config.Load()
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError && cfg == nil {
				t.Error("expected valid config, got nil")
			}
		})
	}
}

// TestConfigValidation tests configuration validation logic
func TestConfigValidation(t *testing.T) {
	t.Run("sqlite_creates_default_path", func(t *testing.T) {
		os.Setenv("DB_TYPE", "sqlite")
		defer os.Unsetenv("DB_TYPE")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.Database.Path == "" {
			t.Error("expected default DB path to be set")
		}

		if !strings.Contains(cfg.Database.Path, ".task-manager") {
			t.Errorf("expected default path to contain .task-manager, got: %s", cfg.Database.Path)
		}
	})

	t.Run("postgres_requires_credentials", func(t *testing.T) {
		os.Setenv("DB_TYPE", "postgres")
		os.Setenv("DB_HOST", "localhost")
		defer func() {
			os.Unsetenv("DB_TYPE")
			os.Unsetenv("DB_HOST")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("expected error for missing postgres credentials")
		}

		if !strings.Contains(err.Error(), "database name") && !strings.Contains(err.Error(), "database user") {
			t.Errorf("expected error about missing credentials, got: %v", err)
		}
	})

	t.Run("invalid_log_level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "invalid")
		defer os.Unsetenv("LOG_LEVEL")

		_, err := config.Load()
		if err == nil {
			t.Error("expected error for invalid log level")
		}

		if !strings.Contains(err.Error(), "log level") {
			t.Errorf("expected error about invalid log level, got: %v", err)
		}
	})

	t.Run("invalid_log_format", func(t *testing.T) {
		os.Setenv("LOG_FORMAT", "invalid")
		defer os.Unsetenv("LOG_FORMAT")

		_, err := config.Load()
		if err == nil {
			t.Error("expected error for invalid log format")
		}

		if !strings.Contains(err.Error(), "log format") {
			t.Errorf("expected error about invalid log format, got: %v", err)
		}
	})
}

// TestConfigFileLoading tests YAML configuration file loading
func TestConfigFileLoading(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `database:
  type: sqlite
  path: /tmp/test-config.db

logging:
  level: debug
  format: json
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	os.Setenv("CONFIG_FILE", configPath)
	defer os.Unsetenv("CONFIG_FILE")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config from file: %v", err)
	}

	if cfg.Database.Type != "sqlite" {
		t.Errorf("expected database type sqlite, got %s", cfg.Database.Type)
	}

	if cfg.Database.Path != "/tmp/test-config.db" {
		t.Errorf("expected database path /tmp/test-config.db, got %s", cfg.Database.Path)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("expected log level debug, got %s", cfg.Logging.Level)
	}

	if cfg.Logging.Format != "json" {
		t.Errorf("expected log format json, got %s", cfg.Logging.Format)
	}
}

// TestEnvVarOverride tests that environment variables override config file
func TestEnvVarOverride(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `database:
  type: sqlite
  path: /tmp/config-file.db

logging:
  level: info
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	os.Setenv("CONFIG_FILE", configPath)
	os.Setenv("DB_PATH", "/tmp/env-override.db")
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("DB_PATH")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Environment variables should override config file
	if cfg.Database.Path != "/tmp/env-override.db" {
		t.Errorf("expected env var to override config file for DB path, got %s", cfg.Database.Path)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("expected env var to override config file for log level, got %s", cfg.Logging.Level)
	}
}

// TestConfigMissingFile tests handling of missing config file
func TestConfigMissingFile(t *testing.T) {
	os.Setenv("CONFIG_FILE", "/nonexistent/config.yaml")
	defer os.Unsetenv("CONFIG_FILE")

	_, err := config.Load()
	if err == nil {
		t.Error("expected error for explicitly specified missing config file")
	}
}

// TestConfigOptionalFile tests that optional config file doesn't cause error
func TestConfigOptionalFile(t *testing.T) {
	// Don't set CONFIG_FILE env var, but reference a non-existent default
	// This should not fail as config file is optional when not explicitly specified
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

	cfg, err := config.Load()
	if err != nil {
		t.Errorf("config load should succeed with defaults when optional config file missing: %v", err)
	}

	if cfg == nil {
		t.Error("expected valid config with defaults")
	}
}
