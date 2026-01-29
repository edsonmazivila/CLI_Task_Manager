package integration

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/edson-mazvila/task-manager/internal/domain"
	"github.com/edson-mazvila/task-manager/internal/repository"
	"github.com/edson-mazvila/task-manager/internal/service"
	"github.com/edson-mazvila/task-manager/internal/storage"
)

// TestEnvironment holds test infrastructure
type TestEnvironment struct {
	DBPath  string
	Storage *storage.SQLiteStorage
	Repo    domain.TaskRepository
	Service *service.TaskService
	Logger  *slog.Logger
	ctx     context.Context
}

// setupTestEnvironment creates a clean test environment
func setupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_tasks.db")

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Reduce noise in tests
	}))

	ctx := context.Background()

	// Initialize storage
	store, err := storage.NewSQLiteStorage(ctx, dbPath, logger)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}

	// Initialize repository
	repo := repository.NewSQLiteTaskRepository(store.DB(), logger)

	// Initialize service
	svc := service.NewTaskService(repo, logger)

	return &TestEnvironment{
		DBPath:  dbPath,
		Storage: store,
		Repo:    repo,
		Service: svc,
		Logger:  logger,
		ctx:     ctx,
	}
}

// cleanup closes resources
func (te *TestEnvironment) cleanup(t *testing.T) {
	t.Helper()
	if err := te.Storage.Close(); err != nil {
		t.Errorf("failed to close storage: %v", err)
	}
}

// TestTaskLifecycle tests the complete lifecycle of a task
func TestTaskLifecycle(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	// Create a task
	task, err := env.Service.CreateTask(env.ctx, "Integration Test Task", "Test description", domain.TaskPriorityHigh)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	if task.ID == "" {
		t.Error("task ID should not be empty")
	}
	if task.Title != "Integration Test Task" {
		t.Errorf("expected title 'Integration Test Task', got '%s'", task.Title)
	}
	if task.Status != domain.TaskStatusPending {
		t.Errorf("expected status pending, got %s", task.Status)
	}
	if task.Priority != domain.TaskPriorityHigh {
		t.Errorf("expected priority high, got %s", task.Priority)
	}

	// Retrieve the task
	retrieved, err := env.Service.GetTask(env.ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrieved.ID != task.ID {
		t.Errorf("expected ID %s, got %s", task.ID, retrieved.ID)
	}

	// Update the task
	updated, err := env.Service.UpdateTask(env.ctx, task.ID, "Updated Title", "Updated description", domain.TaskPriorityMedium)
	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("expected updated title, got '%s'", updated.Title)
	}
	if updated.Priority != domain.TaskPriorityMedium {
		t.Errorf("expected priority medium, got %s", updated.Priority)
	}

	// Complete the task
	completed, err := env.Service.CompleteTask(env.ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}

	if completed.Status != domain.TaskStatusCompleted {
		t.Errorf("expected status completed, got %s", completed.Status)
	}
	if completed.CompletedAt == nil {
		t.Error("completed_at should be set")
	}

	// Delete the task
	err = env.Service.DeleteTask(env.ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	// Verify deletion
	_, err = env.Service.GetTask(env.ctx, task.ID)
	if err != domain.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

// TestTaskFiltering tests various filtering scenarios
func TestTaskFiltering(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	// Create multiple tasks
	tasks := []struct {
		title    string
		priority domain.TaskPriority
	}{
		{"High Priority Task 1", domain.TaskPriorityHigh},
		{"High Priority Task 2", domain.TaskPriorityHigh},
		{"Medium Priority Task", domain.TaskPriorityMedium},
		{"Low Priority Task", domain.TaskPriorityLow},
	}

	var createdIDs []string
	for _, tc := range tasks {
		task, err := env.Service.CreateTask(env.ctx, tc.title, "", tc.priority)
		if err != nil {
			t.Fatalf("failed to create task: %v", err)
		}
		createdIDs = append(createdIDs, task.ID)
	}

	// Complete one task
	_, err := env.Service.CompleteTask(env.ctx, createdIDs[0])
	if err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}

	t.Run("filter_by_status_pending", func(t *testing.T) {
		status := domain.TaskStatusPending
		filter := domain.TaskFilter{Status: &status}

		results, err := env.Service.ListTasks(env.ctx, filter)
		if err != nil {
			t.Fatalf("failed to list tasks: %v", err)
		}

		if len(results) != 3 {
			t.Errorf("expected 3 pending tasks, got %d", len(results))
		}

		for _, task := range results {
			if task.Status != domain.TaskStatusPending {
				t.Errorf("expected all tasks to be pending, got %s", task.Status)
			}
		}
	})

	t.Run("filter_by_status_completed", func(t *testing.T) {
		status := domain.TaskStatusCompleted
		filter := domain.TaskFilter{Status: &status}

		results, err := env.Service.ListTasks(env.ctx, filter)
		if err != nil {
			t.Fatalf("failed to list tasks: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("expected 1 completed task, got %d", len(results))
		}

		if len(results) > 0 && results[0].Status != domain.TaskStatusCompleted {
			t.Errorf("expected completed status, got %s", results[0].Status)
		}
	})

	t.Run("filter_by_priority_high", func(t *testing.T) {
		priority := domain.TaskPriorityHigh
		filter := domain.TaskFilter{Priority: &priority}

		results, err := env.Service.ListTasks(env.ctx, filter)
		if err != nil {
			t.Fatalf("failed to list tasks: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 high priority tasks, got %d", len(results))
		}

		for _, task := range results {
			if task.Priority != domain.TaskPriorityHigh {
				t.Errorf("expected all tasks to be high priority, got %s", task.Priority)
			}
		}
	})

	t.Run("filter_by_date_range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		tomorrow := now.Add(24 * time.Hour)

		filter := domain.TaskFilter{
			FromDate: &yesterday,
			ToDate:   &tomorrow,
		}

		results, err := env.Service.ListTasks(env.ctx, filter)
		if err != nil {
			t.Fatalf("failed to list tasks: %v", err)
		}

		if len(results) != 4 {
			t.Errorf("expected 4 tasks within date range, got %d", len(results))
		}
	})

	t.Run("combined_filters", func(t *testing.T) {
		status := domain.TaskStatusPending
		priority := domain.TaskPriorityHigh
		filter := domain.TaskFilter{
			Status:   &status,
			Priority: &priority,
		}

		results, err := env.Service.ListTasks(env.ctx, filter)
		if err != nil {
			t.Fatalf("failed to list tasks: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("expected 1 task matching both filters, got %d", len(results))
		}
	})
}

// TestConcurrentOperations tests thread safety
func TestConcurrentOperations(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	const numTasks = 10
	errChan := make(chan error, numTasks)

	// Create tasks concurrently
	for i := 0; i < numTasks; i++ {
		go func(index int) {
			title := fmt.Sprintf("Concurrent Task %d", index)
			_, err := env.Service.CreateTask(env.ctx, title, "", domain.TaskPriorityMedium)
			errChan <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numTasks; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("concurrent task creation failed: %v", err)
		}
	}

	// Verify all tasks were created
	filter := domain.TaskFilter{}
	tasks, err := env.Service.ListTasks(env.ctx, filter)
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}

	if len(tasks) != numTasks {
		t.Errorf("expected %d tasks, got %d", numTasks, len(tasks))
	}
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	t.Run("get_nonexistent_task", func(t *testing.T) {
		_, err := env.Service.GetTask(env.ctx, "nonexistent-id")
		if err != domain.ErrTaskNotFound {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("update_nonexistent_task", func(t *testing.T) {
		_, err := env.Service.UpdateTask(env.ctx, "nonexistent-id", "Title", "", domain.TaskPriorityHigh)
		if err != domain.ErrTaskNotFound {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("complete_nonexistent_task", func(t *testing.T) {
		_, err := env.Service.CompleteTask(env.ctx, "nonexistent-id")
		if err != domain.ErrTaskNotFound {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("delete_nonexistent_task", func(t *testing.T) {
		err := env.Service.DeleteTask(env.ctx, "nonexistent-id")
		if err != domain.ErrTaskNotFound {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("create_task_with_empty_title", func(t *testing.T) {
		_, err := env.Service.CreateTask(env.ctx, "", "Description", domain.TaskPriorityHigh)
		if err == nil {
			t.Error("expected error for empty title, got nil")
		}
	})

	t.Run("create_task_with_invalid_priority", func(t *testing.T) {
		_, err := env.Service.CreateTask(env.ctx, "Title", "", "invalid")
		if err == nil {
			t.Error("expected error for invalid priority, got nil")
		}
	})

	t.Run("invalid_task_id", func(t *testing.T) {
		_, err := env.Service.GetTask(env.ctx, "")
		if err != domain.ErrInvalidTaskID {
			t.Errorf("expected ErrInvalidTaskID, got %v", err)
		}
	})
}

// TestDatabasePersistence tests data persistence across connections
func TestDatabasePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persistence_test.db")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	ctx := context.Background()

	// Create first environment and add tasks
	store1, err := storage.NewSQLiteStorage(ctx, dbPath, logger)
	if err != nil {
		t.Fatalf("failed to initialize first storage: %v", err)
	}

	repo1 := repository.NewSQLiteTaskRepository(store1.DB(), logger)
	svc1 := service.NewTaskService(repo1, logger)

	task, err := svc1.CreateTask(ctx, "Persistent Task", "Should survive", domain.TaskPriorityHigh)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	taskID := task.ID

	// Close first connection
	if err := store1.Close(); err != nil {
		t.Fatalf("failed to close first storage: %v", err)
	}

	// Open new connection and verify data
	store2, err := storage.NewSQLiteStorage(ctx, dbPath, logger)
	if err != nil {
		t.Fatalf("failed to initialize second storage: %v", err)
	}
	defer store2.Close()

	repo2 := repository.NewSQLiteTaskRepository(store2.DB(), logger)
	svc2 := service.NewTaskService(repo2, logger)

	retrieved, err := svc2.GetTask(ctx, taskID)
	if err != nil {
		t.Fatalf("failed to retrieve task after reconnect: %v", err)
	}

	if retrieved.ID != taskID {
		t.Errorf("expected ID %s, got %s", taskID, retrieved.ID)
	}
	if retrieved.Title != "Persistent Task" {
		t.Errorf("expected title 'Persistent Task', got '%s'", retrieved.Title)
	}
}

// TestMigrationIdempotency tests that migrations can be run multiple times safely
func TestMigrationIdempotency(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "migration_test.db")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	ctx := context.Background()

	// Run migrations first time
	store1, err := storage.NewSQLiteStorage(ctx, dbPath, logger)
	if err != nil {
		t.Fatalf("failed to initialize storage first time: %v", err)
	}
	store1.Close()

	// Run migrations second time (should be idempotent)
	store2, err := storage.NewSQLiteStorage(ctx, dbPath, logger)
	if err != nil {
		t.Fatalf("failed to initialize storage second time: %v", err)
	}
	defer store2.Close()

	// Verify database structure
	var tableExists int
	err = store2.DB().QueryRowContext(ctx,
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='tasks'",
	).Scan(&tableExists)

	if err != nil {
		t.Fatalf("failed to check table existence: %v", err)
	}

	if tableExists != 1 {
		t.Errorf("expected tasks table to exist once, got %d", tableExists)
	}
}

// TestTransactionRollback tests transaction rollback on errors
func TestTransactionRollback(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	// Create a valid task
	task, err := env.Service.CreateTask(env.ctx, "Test Task", "", domain.TaskPriorityMedium)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	// Attempt to update with invalid data should not corrupt the database
	originalTitle := task.Title
	_, err = env.Service.UpdateTask(env.ctx, task.ID, "", "", "invalid-priority")
	if err == nil {
		t.Error("expected error for invalid update, got nil")
	}

	// Verify original data is intact
	retrieved, err := env.Service.GetTask(env.ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrieved.Title != originalTitle {
		t.Errorf("task was corrupted: expected title '%s', got '%s'", originalTitle, retrieved.Title)
	}
}

// TestTaskTimestamps tests that timestamps are set correctly
func TestTaskTimestamps(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	beforeCreate := time.Now()

	task, err := env.Service.CreateTask(env.ctx, "Timestamp Test", "", domain.TaskPriorityMedium)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	afterCreate := time.Now()

	// Check creation timestamp
	if task.CreatedAt.Before(beforeCreate) || task.CreatedAt.After(afterCreate) {
		t.Error("created_at timestamp is outside expected range")
	}

	if task.UpdatedAt.Before(beforeCreate) || task.UpdatedAt.After(afterCreate) {
		t.Error("updated_at timestamp is outside expected range")
	}

	if task.CompletedAt != nil {
		t.Error("completed_at should be nil for new task")
	}

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Update the task
	beforeUpdate := time.Now()
	updated, err := env.Service.UpdateTask(env.ctx, task.ID, "Updated Title", "", domain.TaskPriorityHigh)
	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}
	afterUpdate := time.Now()

	if !updated.UpdatedAt.After(task.UpdatedAt) {
		t.Error("updated_at should be newer after update")
	}

	if updated.UpdatedAt.Before(beforeUpdate) || updated.UpdatedAt.After(afterUpdate) {
		t.Error("updated_at timestamp is outside expected range after update")
	}

	// Complete the task
	beforeComplete := time.Now()
	completed, err := env.Service.CompleteTask(env.ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}
	afterComplete := time.Now()

	if completed.CompletedAt == nil {
		t.Fatal("completed_at should be set after completion")
	}

	if completed.CompletedAt.Before(beforeComplete) || completed.CompletedAt.After(afterComplete) {
		t.Error("completed_at timestamp is outside expected range")
	}

	if !completed.UpdatedAt.After(updated.UpdatedAt) {
		t.Error("updated_at should be newer after completion")
	}
}

// TestListTasksOrdering tests that tasks are ordered correctly
func TestListTasksOrdering(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	// Create tasks with delays to ensure different timestamps
	for i := 1; i <= 3; i++ {
		title := fmt.Sprintf("Task %d", i)
		_, err := env.Service.CreateTask(env.ctx, title, "", domain.TaskPriorityMedium)
		if err != nil {
			t.Fatalf("failed to create task: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	filter := domain.TaskFilter{}
	tasks, err := env.Service.ListTasks(env.ctx, filter)
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	// Verify descending order (newest first)
	for i := 0; i < len(tasks)-1; i++ {
		if tasks[i].CreatedAt.Before(tasks[i+1].CreatedAt) {
			t.Error("tasks should be ordered by created_at DESC")
		}
	}
}

// TestDatabaseIndexes tests that indexes are created
func TestDatabaseIndexes(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.cleanup(t)

	// Query for indexes
	rows, err := env.Storage.DB().Query(
		"SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='tasks'",
	)
	if err != nil {
		t.Fatalf("failed to query indexes: %v", err)
	}
	defer rows.Close()

	indexes := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("failed to scan index name: %v", err)
		}
		indexes[name] = true
	}

	expectedIndexes := []string{
		"idx_tasks_status",
		"idx_tasks_priority",
		"idx_tasks_created_at",
	}

	for _, idx := range expectedIndexes {
		if !indexes[idx] {
			t.Errorf("expected index %s to exist", idx)
		}
	}
}

// BenchmarkTaskCreation benchmarks task creation performance
func BenchmarkTaskCreation(b *testing.B) {
	env := setupTestEnvironment(&testing.T{})
	defer env.cleanup(&testing.T{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		title := fmt.Sprintf("Benchmark Task %d", i)
		_, err := env.Service.CreateTask(env.ctx, title, "", domain.TaskPriorityMedium)
		if err != nil {
			b.Fatalf("failed to create task: %v", err)
		}
	}
}

// BenchmarkTaskQuery benchmarks task retrieval performance
func BenchmarkTaskQuery(b *testing.B) {
	env := setupTestEnvironment(&testing.T{})
	defer env.cleanup(&testing.T{})

	// Create a task to query
	task, err := env.Service.CreateTask(env.ctx, "Benchmark Task", "", domain.TaskPriorityMedium)
	if err != nil {
		b.Fatalf("failed to create task: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := env.Service.GetTask(env.ctx, task.ID)
		if err != nil {
			b.Fatalf("failed to get task: %v", err)
		}
	}
}

// BenchmarkListTasks benchmarks listing performance
func BenchmarkListTasks(b *testing.B) {
	env := setupTestEnvironment(&testing.T{})
	defer env.cleanup(&testing.T{})

	// Create 100 tasks
	for i := 0; i < 100; i++ {
		title := fmt.Sprintf("Task %d", i)
		_, err := env.Service.CreateTask(env.ctx, title, "", domain.TaskPriorityMedium)
		if err != nil {
			b.Fatalf("failed to create task: %v", err)
		}
	}

	filter := domain.TaskFilter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := env.Service.ListTasks(env.ctx, filter)
		if err != nil {
			b.Fatalf("failed to list tasks: %v", err)
		}
	}
}
