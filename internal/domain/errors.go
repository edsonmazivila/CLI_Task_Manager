package domain

import "errors"

var (
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")

	// ErrInvalidTaskID is returned when a task ID is invalid
	ErrInvalidTaskID = errors.New("invalid task ID")

	// ErrDuplicateTask is returned when trying to create a duplicate task
	ErrDuplicateTask = errors.New("duplicate task")
)
