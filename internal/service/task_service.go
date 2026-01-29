package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/edson-mazvila/task-manager/internal/domain"
	"github.com/google/uuid"
)

// TaskService provides business logic for task management
type TaskService struct {
	repo   domain.TaskRepository
	logger *slog.Logger
}

// NewTaskService creates a new task service
func NewTaskService(repo domain.TaskRepository, logger *slog.Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, title, description string, priority domain.TaskPriority) (*domain.Task, error) {
	task := &domain.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      domain.TaskStatusPending,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := task.Validate(); err != nil {
		s.logger.Warn("Task validation failed", "error", err)
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	if err := s.repo.Create(ctx, task); err != nil {
		s.logger.Error("Failed to create task", "error", err)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	s.logger.Info("Task created successfully", "task_id", task.ID, "title", task.Title)
	return task, nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	if id == "" {
		return nil, domain.ErrInvalidTaskID
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get task", "error", err, "task_id", id)
		return nil, err
	}

	return task, nil
}

// ListTasks retrieves all tasks based on filter criteria
func (s *TaskService) ListTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error) {
	tasks, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list tasks", "error", err)
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	s.logger.Debug("Tasks listed", "count", len(tasks))
	return tasks, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(ctx context.Context, id, title, description string, priority domain.TaskPriority) (*domain.Task, error) {
	if id == "" {
		return nil, domain.ErrInvalidTaskID
	}

	// Get existing task
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get task for update", "error", err, "task_id", id)
		return nil, err
	}

	// Update fields
	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if priority != "" {
		task.Priority = priority
	}
	task.UpdatedAt = time.Now()

	// Validate updated task
	if err := task.Validate(); err != nil {
		s.logger.Warn("Task validation failed", "error", err)
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	// Save updated task
	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("Failed to update task", "error", err, "task_id", id)
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	s.logger.Info("Task updated successfully", "task_id", task.ID)
	return task, nil
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(ctx context.Context, id string) (*domain.Task, error) {
	if id == "" {
		return nil, domain.ErrInvalidTaskID
	}

	// Get existing task
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get task for completion", "error", err, "task_id", id)
		return nil, err
	}

	// Check if already completed
	if task.Status == domain.TaskStatusCompleted {
		s.logger.Warn("Task already completed", "task_id", id)
		return task, nil
	}

	// Mark as completed
	task.MarkCompleted()

	// Save updated task
	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("Failed to complete task", "error", err, "task_id", id)
		return nil, fmt.Errorf("failed to complete task: %w", err)
	}

	s.logger.Info("Task completed successfully", "task_id", task.ID)
	return task, nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidTaskID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete task", "error", err, "task_id", id)
		return err
	}

	s.logger.Info("Task deleted successfully", "task_id", id)
	return nil
}
