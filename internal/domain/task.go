package domain

import (
	"context"
	"errors"
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusCompleted TaskStatus = "completed"
)

// TaskPriority represents the priority level of a task
type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

// Task represents a task in the system
type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

// TaskFilter contains filter criteria for querying tasks
type TaskFilter struct {
	Status   *TaskStatus
	Priority *TaskPriority
	FromDate *time.Time
	ToDate   *time.Time
}

// Validate validates the task
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("task title cannot be empty")
	}

	if t.Status != TaskStatusPending && t.Status != TaskStatusCompleted {
		return errors.New("invalid task status")
	}

	if t.Priority != TaskPriorityLow && t.Priority != TaskPriorityMedium && t.Priority != TaskPriorityHigh {
		return errors.New("invalid task priority")
	}

	return nil
}

// MarkCompleted marks the task as completed
func (t *Task) MarkCompleted() {
	t.Status = TaskStatusCompleted
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context, filter TaskFilter) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}
