package memory

import (
	"context"
	"task-tracking-service/internal/core/domain"
	"task-tracking-service/pkg/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_Create(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	t.Run("successfully creates task", func(t *testing.T) {
		// Create a test task
		task := &domain.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      domain.StatusPending,
			DueDate:     time.Now().Add(24 * time.Hour),
		}

		// Create the task
		err := repo.Create(ctx, task)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, task.ID, "Task ID should be generated")

		// Verify task was stored
		storedTask, err := repo.GetByID(ctx, task.ID)
		assert.NoError(t, err)
		assert.NotNil(t, storedTask)
		assert.Equal(t, task.Title, storedTask.Title)
		assert.Equal(t, task.Description, storedTask.Description)
		assert.Equal(t, task.Status, storedTask.Status)
	})

	t.Run("creates unique IDs for different tasks", func(t *testing.T) {
		// Create two tasks
		task1 := &domain.Task{
			Title: "Task 1",
		}
		task2 := &domain.Task{
			Title: "Task 2",
		}

		// Create both tasks
		err1 := repo.Create(ctx, task1)
		err2 := repo.Create(ctx, task2)

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, task1.ID)
		assert.NotEmpty(t, task2.ID)
		assert.NotEqual(t, task1.ID, task2.ID, "Task IDs should be unique")
	})

	t.Run("creates independent copies of tasks", func(t *testing.T) {
		// Create initial task
		originalTask := &domain.Task{
			Title: "Original Task",
		}
		err := repo.Create(ctx, originalTask)
		assert.NoError(t, err)

		// Modify the original task after creation
		originalTask.Title = "Modified Title"

		// Get the stored task
		storedTask, err := repo.GetByID(ctx, originalTask.ID)
		assert.NoError(t, err)

		// Verify the stored task wasn't modified
		assert.Equal(t, "Original Task", storedTask.Title)
	})
}

func TestTaskRepository_GetByID(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	tests := []struct {
		name          string
		setupTask     *domain.Task
		taskID        string
		expectedError error
	}{
		{
			name: "successfully gets existing task",
			setupTask: &domain.Task{
				Title:       "Test Task",
				Description: "Test Description",
				Status:      domain.StatusPending,
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			taskID:        "", // Will be set after creation
			expectedError: nil,
		},
		{
			name:          "returns error for non-existent task",
			setupTask:     nil,
			taskID:        "non-existent-id",
			expectedError: &errors.NotFoundError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTask != nil {
				err := repo.Create(ctx, tt.setupTask)
				assert.NoError(t, err)
				tt.taskID = tt.setupTask.ID
			}

			task, err := repo.GetByID(ctx, tt.taskID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tt.setupTask.ID, task.ID)
				assert.Equal(t, tt.setupTask.Title, task.Title)
				assert.Equal(t, tt.setupTask.Description, task.Description)
			}
		})
	}
}

func TestTaskRepository_List(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	t.Run("returns empty list when no tasks exist", func(t *testing.T) {
		tasks, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("returns all created tasks", func(t *testing.T) {
		// Create multiple tasks
		task1 := &domain.Task{Title: "Task 1", Status: domain.StatusPending}
		task2 := &domain.Task{Title: "Task 2", Status: domain.StatusPending}

		err := repo.Create(ctx, task1)
		assert.NoError(t, err)
		err = repo.Create(ctx, task2)
		assert.NoError(t, err)

		tasks, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
	})
}

func TestTaskRepository_Update(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	t.Run("successfully updates existing task", func(t *testing.T) {
		// Create initial task
		task := &domain.Task{
			Title:       "Original Title",
			Description: "Original Description",
			Status:      domain.StatusPending,
			DueDate:     time.Now().Add(24 * time.Hour),
		}
		err := repo.Create(ctx, task)
		assert.NoError(t, err)

		// Update task
		updatedTask := &domain.Task{
			ID:          task.ID,
			Title:       "Updated Title",
			Description: "Updated Description",
			Status:      domain.StatusInProgress,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		err = repo.Update(ctx, updatedTask)
		assert.NoError(t, err)

		// Retrieve updated task
		retrieved, err := repo.GetByID(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", retrieved.Title)
		assert.Equal(t, "Updated Description", retrieved.Description)
		assert.Equal(t, domain.StatusInProgress, retrieved.Status)
	})

	t.Run("fails to update non-existent task", func(t *testing.T) {
		nonExistentTask := &domain.Task{
			ID:    "non-existent-id",
			Title: "Non-existent Task",
		}

		err := repo.Update(ctx, nonExistentTask)
		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err))
	})

	t.Run("updates maintain data independence", func(t *testing.T) {
		// Create initial task
		task := &domain.Task{
			Title: "Original Task",
		}
		err := repo.Create(ctx, task)
		assert.NoError(t, err)

		// Create update with modifications
		updateTask := &domain.Task{
			ID:    task.ID,
			Title: "Updated Task",
		}
		err = repo.Update(ctx, updateTask)
		assert.NoError(t, err)

		// Modify the update task after updating
		updateTask.Title = "Modified After Update"

		// Verify stored task wasn't affected by post-update modification
		stored, err := repo.GetByID(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", stored.Title)
	})
}

func TestTaskRepository_Delete(t *testing.T) {
	repo := NewTaskRepository()
	ctx := context.Background()

	t.Run("successfully deletes existing task", func(t *testing.T) {
		task := &domain.Task{Title: "Task to Delete"}
		err := repo.Create(ctx, task)
		assert.NoError(t, err)

		err = repo.Delete(ctx, task.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(ctx, task.ID)
		assert.Error(t, err)
		assert.IsType(t, &errors.NotFoundError{}, err)
	})

	t.Run("returns error when deleting non-existent task", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.IsType(t, &errors.NotFoundError{}, err)
	})
}
