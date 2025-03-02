package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"task-tracking-service/internal/config"
	"task-tracking-service/internal/core/domain"
	customerrors "task-tracking-service/pkg/errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Set GO_ENV to test to skip loading .env file
	os.Setenv("GO_ENV", "test")

	// Get current user to use as fallback
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		// If USER env var isn't set, try to get username another way
		currentUser = "postgres" // fallback
	}

	// Set test environment variables if not set
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", currentUser) // Use current user as default
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "dummy_password_for_testing") // Dummy password to satisfy validation
	}
	if os.Getenv("DB_SSL_MODE") == "" {
		os.Setenv("DB_SSL_MODE", "disable")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "taskdb_test")
	}

	// Set required config values to pass validation
	os.Setenv("APP_ENV", "development")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_READ_TIMEOUT", "60s")
	os.Setenv("SERVER_WRITE_TIMEOUT", "60s")
	os.Setenv("SERVER_BASE_URL", "http://localhost:8080")
	os.Setenv("API_BASE_PATH", "/api")
	os.Setenv("API_KEY", "test-api-key-at-least-32-characters-long")
	os.Setenv("CORS_ALLOWED_ORIGINS", "*")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOG_FORMAT", "text")
	os.Setenv("DB_MAX_OPEN_CONNS", "10")
	os.Setenv("DB_MAX_IDLE_CONNS", "5")
	os.Setenv("DB_CONN_MAX_LIFETIME", "5m")
	os.Setenv("REPOSITORY_TYPE", "postgres")

	// Load config
	cfg, err := config.Load()
	require.NoError(t, err)

	// Try direct connection to test database first
	testDBName := "taskdb_test"

	testConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(string(cfg.Database.User)),
		url.QueryEscape(string(cfg.Database.Password)),
		cfg.Database.Host,
		cfg.Database.Port,
		testDBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", testConnString)
	if err == nil {
		// If we can connect to the test database, just use it
		err = db.Ping()
		if err == nil {
			// Successfully connected to existing test database
			setupTestTables(t, db)
			return db
		}
		db.Close()
	}

	// If direct connection fails, try to create the database
	// Connect to default postgres database
	defaultConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		url.QueryEscape(string(cfg.Database.User)),
		url.QueryEscape(string(cfg.Database.Password)),
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	defaultDB, err := sql.Open("postgres", defaultConnString)
	if err != nil {
		// Try connecting to another default database
		defaultConnString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/template1?sslmode=%s",
			url.QueryEscape(string(cfg.Database.User)),
			url.QueryEscape(string(cfg.Database.Password)),
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
		defaultDB, err = sql.Open("postgres", defaultConnString)
		require.NoError(t, err, "Failed to connect to any default database")
	}

	defer defaultDB.Close()
	require.NoError(t, defaultDB.Ping(), "Failed to ping default database")

	// Create test database if it doesn't exist (without specifying owner)
	_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		// If error contains "already exists", it's fine; otherwise, fail the test
		if !strings.Contains(err.Error(), "already exists") {
			require.NoError(t, err, "Failed to create test database")
		}
	}

	// Now connect to the test database after creating it
	db, err = sql.Open("postgres", testConnString)
	require.NoError(t, err, "Failed to connect to test database")
	require.NoError(t, db.Ping(), "Failed to ping test database")

	// Set up the test tables
	setupTestTables(t, db)

	return db
}

// setupTestTables creates and truncates tables for testing
func setupTestTables(t *testing.T, db *sql.DB) {
	// Create tasks table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			due_date TIMESTAMP
		)
	`)
	require.NoError(t, err, "Failed to create tasks table")

	// Clear the tasks table for a fresh test
	_, err = db.Exec("TRUNCATE TABLE tasks")
	require.NoError(t, err, "Failed to truncate tasks table")
}

func TestTaskRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTaskRepository(db)
	ctx := context.Background()

	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusPending,
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "Test Description", task.Description)
	assert.Equal(t, domain.StatusPending, task.Status)
	assert.NotEmpty(t, task.CreatedAt)
	assert.NotEmpty(t, task.UpdatedAt)

	// Test getting the task
	fetchedTask, err := repo.GetByID(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.ID, fetchedTask.ID)
	assert.Equal(t, task.Title, fetchedTask.Title)

	// Test getting non-existent task
	_, err = repo.GetByID(ctx, uuid.New().String())
	assert.ErrorIs(t, err, customerrors.ErrTaskNotFound)
}

func TestTaskRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTaskRepository(db)
	ctx := context.Background()

	// Create a task first
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusPending,
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	// Test getting the task
	fetchedTask, err := repo.GetByID(ctx, task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.ID, fetchedTask.ID)
	assert.Equal(t, task.Title, fetchedTask.Title)

	// Test getting non-existent task
	_, err = repo.GetByID(ctx, uuid.New().String())
	assert.ErrorIs(t, err, customerrors.ErrTaskNotFound)
}

// Add more tests for List, Update, and Delete...
