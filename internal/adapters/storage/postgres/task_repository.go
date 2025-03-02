package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"your-module/internal/core/domain"
	"your-module/internal/core/ports"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// Create stores a new task in the database
func (r *TaskRepository) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	query := `
		INSERT INTO tasks (id, title, description, status, created_at, updated_at, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, status, created_at, updated_at, due_date`

	row := r.db.QueryRowContext(
		ctx,
		query,
		uuid.New(),
		task.Title,
		task.Description,
		task.Status,
		time.Now(),
		time.Now(),
		task.DueDate,
	)

	var savedTask domain.Task
	err := row.Scan(
		&savedTask.ID,
		&savedTask.Title,
		&savedTask.Description,
		&savedTask.Status,
		&savedTask.CreatedAt,
		&savedTask.UpdatedAt,
		&savedTask.DueDate,
	)

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return savedTask, nil
}

// Get retrieves a task by ID from the database
func (r *TaskRepository) Get(ctx context.Context, id string) (domain.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at, due_date
		FROM tasks
		WHERE id = $1`

	var task domain.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DueDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, ports.ErrTaskNotFound
		}
		return domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// List retrieves all tasks from the database
func (r *TaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at, due_date
		FROM tasks
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update modifies an existing task in the database
func (r *TaskRepository) Update(ctx context.Context, task domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4, due_date = $5
		WHERE id = $6`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		time.Now(),
		task.DueDate,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrTaskNotFound
	}

	return nil
}

// Delete removes a task from the database
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrTaskNotFound
	}

	return nil
}
