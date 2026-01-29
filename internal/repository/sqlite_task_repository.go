// Package repository provides data access implementations for the task manager.
// This layer handles all database interactions and SQL query execution,
// following the repository pattern to abstract persistence details from business logic.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/edson-mazvila/task-manager/internal/domain"
)

// SQLiteTaskRepository implements TaskRepository interface for SQLite database.
// It provides CRUD operations with transaction support and proper error handling.
// All SQL queries use parameterized statements to prevent SQL injection.
type SQLiteTaskRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewSQLiteTaskRepository creates a new SQLite task repository
func NewSQLiteTaskRepository(db *sql.DB, logger *slog.Logger) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new task into the database.
// Uses parameterized queries to prevent SQL injection and ensure data safety.
// All timestamps are stored in UTC format for consistency across time zones.
func (r *SQLiteTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, created_at, updated_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
		task.CompletedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create task", "error", err, "task_id", task.ID)
		return fmt.Errorf("failed to create task: %w", err)
	}

	r.logger.Info("Task created", "task_id", task.ID)
	return nil
}

// GetByID retrieves a task by its ID
func (r *SQLiteTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, created_at, updated_at, completed_at
		FROM tasks
		WHERE id = ?
	`

	task := &domain.Task{}
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		r.logger.Error("Failed to get task", "error", err, "task_id", id)
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return task, nil
}

// List retrieves tasks based on filter criteria
func (r *SQLiteTaskRepository) List(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error) {
	query := "SELECT id, title, description, status, priority, created_at, updated_at, completed_at FROM tasks WHERE 1=1"
	args := []interface{}{}

	if filter.Status != nil {
		query += " AND status = ?"
		args = append(args, *filter.Status)
	}

	if filter.Priority != nil {
		query += " AND priority = ?"
		args = append(args, *filter.Priority)
	}

	if filter.FromDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *filter.FromDate)
	}

	if filter.ToDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *filter.ToDate)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list tasks", "error", err)
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		var completedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
			&completedAt,
		)

		if err != nil {
			r.logger.Error("Failed to scan task", "error", err)
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating tasks", "error", err)
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update updates an existing task
func (r *SQLiteTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	// First check if the task exists
	existing, err := r.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return domain.ErrTaskNotFound
	}

	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?, updated_at = ?, completed_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.UpdatedAt,
		task.CompletedAt,
		task.ID,
	)

	if err != nil {
		r.logger.Error("Failed to update task", "error", err, "task_id", task.ID)
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	r.logger.Info("Task updated", "task_id", task.ID)
	return nil
}

// Delete deletes a task by its ID
func (r *SQLiteTaskRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM tasks WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete task", "error", err, "task_id", id)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	r.logger.Info("Task deleted", "task_id", id)
	return nil
}
