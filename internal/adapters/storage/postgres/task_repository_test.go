package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"your-module/internal/core/domain"
	"your-module/internal/core/ports"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use test database connection string from environment or config
	db, err := sql.Open("postgres", "postgres://localhost:5432/taskdb_test?sslmode=disable")
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	// Clear the tasks table
	_, err = db.Exec("TRUNCATE TABLE tasks")
	require.NoError(t, err)

	return db
}

func TestTaskRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTaskRepository(db)
	ctx := context.Background()

	task := domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.TaskStatusPending,
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	savedTask, err := repo.Create(ctx, task)
	require.NoError(t, err)

	assert.NotEmpty(t, savedTask.ID)
	assert.Equal(t, task.Title, savedTask.Title)
	assert.Equal(t, task.Description, savedTask.Description)
	assert.Equal(t, task.Status, savedTask.Status)
	assert.NotEmpty(t, savedTask.CreatedAt)
	assert.NotEmpty(t, savedTask.UpdatedAt)
	assert.Equal(t, task.DueDate.Unix(), savedTask.DueDate.Unix())
}

func TestTaskRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTaskRepository(db)
	ctx := context.Background()

	// Create a task first
	task := domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.TaskStatusPending,
	}

	savedTask, err := repo.Create(ctx, task)
	require.NoError(t, err)

	// Test getting the task
	fetchedTask, err := repo.Get(ctx, savedTask.ID)
	require.NoError(t, err)
	assert.Equal(t, savedTask.ID, fetchedTask.ID)
	assert.Equal(t, savedTask.Title, fetchedTask.Title)

	// Test getting non-existent task
	_, err = repo.Get(ctx, uuid.New().String())
	assert.ErrorIs(t, err, ports.ErrTaskNotFound)
}

// Add more tests for List, Update, and Delete...
