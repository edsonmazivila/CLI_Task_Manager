// Package cli provides the command-line interface for the task manager.
// It uses the Cobra framework to provide a rich CLI experience with subcommands,
// flags, and formatted output. This is the presentation layer that interacts
// with users and delegates work to the service layer.
package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"
	"time"

	"github.com/edson-mazvila/task-manager/internal/domain"
	"github.com/edson-mazvila/task-manager/internal/service"
	"github.com/spf13/cobra"
)

// CLI holds the CLI configuration and dependencies.
// It follows dependency injection principles, receiving the service layer
// and logger through the constructor to maintain loose coupling.
type CLI struct {
	service *service.TaskService
	logger  *slog.Logger
}

// NewCLI creates a new CLI instance
func NewCLI(service *service.TaskService, logger *slog.Logger) *CLI {
	return &CLI{
		service: service,
		logger:  logger,
	}
}

// RootCmd returns the root command with all subcommands attached.
// Subcommands include: add, list, get, update, complete, delete.
// Each command has its own flags and validation logic.
func (c *CLI) RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "A production-grade CLI task manager",
		Long:  `Task Manager is a CLI application for managing your tasks efficiently.`,
	}

	rootCmd.AddCommand(
		c.addCmd(),
		c.listCmd(),
		c.completeCmd(),
		c.deleteCmd(),
		c.updateCmd(),
		c.getCmd(),
	)

	return rootCmd
}

// addCmd creates the add command
func (c *CLI) addCmd() *cobra.Command {
	var priority string
	var description string

	cmd := &cobra.Command{
		Use:   "add [title]",
		Short: "Add a new task",
		Long:  `Add a new task with the specified title, priority, and optional description.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]

			// Parse priority
			taskPriority := domain.TaskPriority(priority)
			if taskPriority != domain.TaskPriorityLow &&
				taskPriority != domain.TaskPriorityMedium &&
				taskPriority != domain.TaskPriorityHigh {
				return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", priority)
			}

			// Create task
			ctx := context.Background()
			task, err := c.service.CreateTask(ctx, title, description, taskPriority)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			fmt.Printf("✓ Task created successfully\n")
			fmt.Printf("  ID:       %s\n", task.ID)
			fmt.Printf("  Title:    %s\n", task.Title)
			fmt.Printf("  Priority: %s\n", task.Priority)
			if task.Description != "" {
				fmt.Printf("  Description: %s\n", task.Description)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Task priority (low, medium, high)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Task description")

	return cmd
}

// listCmd creates the list command
func (c *CLI) listCmd() *cobra.Command {
	var status string
	var priority string
	var fromDate string
	var toDate string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Long:  `List all tasks with optional filtering by status, priority, and date range.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := domain.TaskFilter{}

			// Parse status filter
			if status != "" {
				taskStatus := domain.TaskStatus(status)
				if taskStatus != domain.TaskStatusPending && taskStatus != domain.TaskStatusCompleted {
					return fmt.Errorf("invalid status: %s (must be pending or completed)", status)
				}
				filter.Status = &taskStatus
			}

			// Parse priority filter
			if priority != "" {
				taskPriority := domain.TaskPriority(priority)
				if taskPriority != domain.TaskPriorityLow &&
					taskPriority != domain.TaskPriorityMedium &&
					taskPriority != domain.TaskPriorityHigh {
					return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", priority)
				}
				filter.Priority = &taskPriority
			}

			// Parse date filters
			if fromDate != "" {
				t, err := time.Parse("2006-01-02", fromDate)
				if err != nil {
					return fmt.Errorf("invalid from-date format (use YYYY-MM-DD): %w", err)
				}
				filter.FromDate = &t
			}

			if toDate != "" {
				t, err := time.Parse("2006-01-02", toDate)
				if err != nil {
					return fmt.Errorf("invalid to-date format (use YYYY-MM-DD): %w", err)
				}
				filter.ToDate = &t
			}

			// List tasks
			ctx := context.Background()
			tasks, err := c.service.ListTasks(ctx, filter)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return nil
			}

			// Display tasks in table format
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPRIORITY\tCREATED")
			fmt.Fprintln(w, "--\t-----\t------\t--------\t-------")

			for _, task := range tasks {
				createdAt := task.CreatedAt.Format("2006-01-02 15:04")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					task.ID[:8], task.Title, task.Status, task.Priority, createdAt)
			}

			w.Flush()
			fmt.Printf("\nTotal: %d task(s)\n", len(tasks))

			return nil
		},
	}

	cmd.Flags().StringVarP(&status, "status", "s", "", "Filter by status (pending, completed)")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "Filter by priority (low, medium, high)")
	cmd.Flags().StringVar(&fromDate, "from", "", "Filter by from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&toDate, "to", "", "Filter by to date (YYYY-MM-DD)")

	return cmd
}

// getCmd creates the get command
func (c *CLI) getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [task-id]",
		Short: "Get task details",
		Long:  `Get detailed information about a specific task.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			ctx := context.Background()
			task, err := c.service.GetTask(ctx, taskID)
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			fmt.Printf("Task Details:\n")
			fmt.Printf("  ID:          %s\n", task.ID)
			fmt.Printf("  Title:       %s\n", task.Title)
			fmt.Printf("  Description: %s\n", task.Description)
			fmt.Printf("  Status:      %s\n", task.Status)
			fmt.Printf("  Priority:    %s\n", task.Priority)
			fmt.Printf("  Created:     %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Updated:     %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))

			if task.CompletedAt != nil {
				fmt.Printf("  Completed:   %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
			}

			return nil
		},
	}

	return cmd
}

// completeCmd creates the complete command
func (c *CLI) completeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete [task-id]",
		Short: "Mark a task as completed",
		Long:  `Mark the specified task as completed.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			ctx := context.Background()
			task, err := c.service.CompleteTask(ctx, taskID)
			if err != nil {
				return fmt.Errorf("failed to complete task: %w", err)
			}

			fmt.Printf("✓ Task marked as completed\n")
			fmt.Printf("  ID:    %s\n", task.ID)
			fmt.Printf("  Title: %s\n", task.Title)

			return nil
		},
	}

	return cmd
}

// deleteCmd creates the delete command
func (c *CLI) deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [task-id]",
		Short: "Delete a task",
		Long:  `Delete the specified task permanently.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			ctx := context.Background()
			if err := c.service.DeleteTask(ctx, taskID); err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			fmt.Printf("✓ Task deleted successfully (ID: %s)\n", taskID)

			return nil
		},
	}

	return cmd
}

// updateCmd creates the update command
func (c *CLI) updateCmd() *cobra.Command {
	var title string
	var description string
	var priority string

	cmd := &cobra.Command{
		Use:   "update [task-id]",
		Short: "Update a task",
		Long:  `Update the specified task's title, description, or priority.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			// At least one field must be provided
			if title == "" && description == "" && priority == "" {
				return fmt.Errorf("at least one field must be provided (--title, --description, or --priority)")
			}

			// Parse priority if provided
			var taskPriority domain.TaskPriority
			if priority != "" {
				taskPriority = domain.TaskPriority(priority)
				if taskPriority != domain.TaskPriorityLow &&
					taskPriority != domain.TaskPriorityMedium &&
					taskPriority != domain.TaskPriorityHigh {
					return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", priority)
				}
			}

			// Update task
			ctx := context.Background()
			task, err := c.service.UpdateTask(ctx, taskID, title, description, taskPriority)
			if err != nil {
				return fmt.Errorf("failed to update task: %w", err)
			}

			fmt.Printf("✓ Task updated successfully\n")
			fmt.Printf("  ID:       %s\n", task.ID)
			fmt.Printf("  Title:    %s\n", task.Title)
			fmt.Printf("  Priority: %s\n", task.Priority)

			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "New task title")
	cmd.Flags().StringVarP(&description, "description", "d", "", "New task description")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "New task priority (low, medium, high)")

	return cmd
}
