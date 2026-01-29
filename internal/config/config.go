package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type     string `yaml:"type"`     // sqlite or postgres
	Path     string `yaml:"path"`     // for SQLite
	Host     string `yaml:"host"`     // for PostgreSQL
	Port     int    `yaml:"port"`     // for PostgreSQL
	Name     string `yaml:"name"`     // for PostgreSQL
	User     string `yaml:"user"`     // for PostgreSQL
	Password string `yaml:"password"` // for PostgreSQL
	SSLMode  string `yaml:"ssl_mode"` // for PostgreSQL
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json or text
}

// Load loads configuration from environment variables and config file
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Type:     getEnvOrDefault("DB_TYPE", "sqlite"),
			Path:     getEnvOrDefault("DB_PATH", ""),
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvIntOrDefault("DB_PORT", 5432),
			Name:     getEnvOrDefault("DB_NAME", "taskmanager"),
			User:     getEnvOrDefault("DB_USER", ""),
			Password: getEnvOrDefault("DB_PASSWORD", ""),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "text"),
		},
	}

	// Store env var overrides before loading config file
	envOverrides := make(map[string]string)
	envVars := []string{"DB_TYPE", "DB_PATH", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSL_MODE", "LOG_LEVEL", "LOG_FORMAT"}
	for _, key := range envVars {
		if val := os.Getenv(key); val != "" {
			envOverrides[key] = val
		}
	}

	// Try to load from config file if it exists
	configPath := getEnvOrDefault("CONFIG_FILE", "config.yaml")
	configExplicit := os.Getenv("CONFIG_FILE") != ""

	if configPath != "" {
		if err := loadFromFile(configPath, cfg, configExplicit); err != nil {
			// Config file is optional, only return error if it was explicitly specified
			if configExplicit {
				return nil, fmt.Errorf("failed to load config file: %w", err)
			}
		}
	}

	// Reapply environment variable overrides
	if _, ok := envOverrides["DB_TYPE"]; ok {
		cfg.Database.Type = envOverrides["DB_TYPE"]
	}
	if _, ok := envOverrides["DB_PATH"]; ok {
		cfg.Database.Path = envOverrides["DB_PATH"]
	}
	if _, ok := envOverrides["DB_HOST"]; ok {
		cfg.Database.Host = envOverrides["DB_HOST"]
	}
	if _, ok := envOverrides["DB_PORT"]; ok {
		cfg.Database.Port = getEnvIntOrDefault("DB_PORT", cfg.Database.Port)
	}
	if _, ok := envOverrides["DB_NAME"]; ok {
		cfg.Database.Name = envOverrides["DB_NAME"]
	}
	if _, ok := envOverrides["DB_USER"]; ok {
		cfg.Database.User = envOverrides["DB_USER"]
	}
	if _, ok := envOverrides["DB_PASSWORD"]; ok {
		cfg.Database.Password = envOverrides["DB_PASSWORD"]
	}
	if _, ok := envOverrides["DB_SSL_MODE"]; ok {
		cfg.Database.SSLMode = envOverrides["DB_SSL_MODE"]
	}
	if _, ok := envOverrides["LOG_LEVEL"]; ok {
		cfg.Logging.Level = envOverrides["LOG_LEVEL"]
	}
	if _, ok := envOverrides["LOG_FORMAT"]; ok {
		cfg.Logging.Format = envOverrides["LOG_FORMAT"]
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(path string, cfg *Config, explicit bool) error {
	if path == "" {
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If file was explicitly specified, return error
		if explicit {
			return fmt.Errorf("config file not found: %s", path)
		}
		return nil // File doesn't exist, skip (implicit path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Type != "sqlite" && c.Database.Type != "postgres" {
		return errors.New("database type must be 'sqlite' or 'postgres'")
	}

	if c.Database.Type == "sqlite" {
		if c.Database.Path == "" {
			// Set default path in user's home directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get user home directory: %w", err)
			}
			c.Database.Path = filepath.Join(homeDir, ".task-manager", "tasks.db")
		}
	}

	if c.Database.Type == "postgres" {
		if c.Database.Host == "" {
			return errors.New("database host is required for PostgreSQL")
		}
		if c.Database.Name == "" {
			return errors.New("database name is required for PostgreSQL")
		}
		if c.Database.User == "" {
			return errors.New("database user is required for PostgreSQL")
		}
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Logging.Level)
	}

	validLogFormats := map[string]bool{"json": true, "text": true}
	if !validLogFormats[c.Logging.Format] {
		return fmt.Errorf("invalid log format: %s (must be json or text)", c.Logging.Format)
	}

	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns the value of an environment variable as int or a default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
