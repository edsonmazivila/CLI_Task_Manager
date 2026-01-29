package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage manages SQLite database connections and migrations
type SQLiteStorage struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(ctx context.Context, dbPath string, logger *slog.Logger) (*SQLiteStorage, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite works best with a single connection
	db.SetMaxIdleConns(1)

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	storage := &SQLiteStorage{
		db:     db,
		logger: logger,
	}

	// Run migrations
	if err := storage.runMigrations(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("SQLite storage initialized", "path", dbPath)

	return storage, nil
}

// DB returns the underlying database connection
func (s *SQLiteStorage) DB() *sql.DB {
	return s.db
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// runMigrations runs database migrations
func (s *SQLiteStorage) runMigrations(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := s.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Inline migration SQL
	migrations := map[string]string{
		"001_create_tasks_table": `
-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed')),
    priority TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    completed_at DATETIME
);

-- Create index on status for faster filtering
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

-- Create index on priority for faster filtering
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);

-- Create index on created_at for faster date filtering
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
		`,
	}

	// Get sorted migration versions
	var versions []string
	for version := range migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	// Apply pending migrations
	for _, version := range versions {

		if appliedMigrations[version] {
			s.logger.Debug("Migration already applied", "version", version)
			continue
		}

		s.logger.Info("Applying migration", "version", version)

		content := migrations[version]

		// Execute migration in a transaction
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute migration SQL
		if _, err := tx.ExecContext(ctx, content); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record migration
		if _, err := tx.ExecContext(ctx, "INSERT INTO migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		s.logger.Info("Migration applied successfully", "version", version)
	}

	return nil
}

// getAppliedMigrations returns a map of applied migration versions
func (s *SQLiteStorage) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT version FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}
