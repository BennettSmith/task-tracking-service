package services

import (
	"context"
	"task-tracking-service/internal/core/domain"
	"task-tracking-service/internal/core/ports"
	"time"
)

type TaskService struct {
	repo ports.TaskRepository
}

func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string, dueDate time.Time) (*domain.Task, error) {
	task := &domain.Task{
		Title:       title,
		Description: description,
		Status:      domain.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     dueDate,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.List(ctx)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	existing, err := s.repo.GetByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	if err := s.validateStatusTransition(existing.Status, task.Status); err != nil {
		return nil, err
	}

	task.CreatedAt = existing.CreatedAt
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *TaskService) validateStatusTransition(from, to domain.TaskStatus) error {
	validTransitions := map[domain.TaskStatus][]domain.TaskStatus{
		domain.StatusPending: {
			domain.StatusInProgress,
			domain.StatusCompleted,
		},
		domain.StatusInProgress: {
			domain.StatusCompleted,
			domain.StatusPending,
		},
		domain.StatusCompleted: {
			domain.StatusPending,
			domain.StatusInProgress,
		},
	}

	if from == to {
		return nil
	}

	allowedTransitions, exists := validTransitions[from]
	if !exists {
		return NewInvalidStatusError("invalid current status")
	}

	for _, allowedStatus := range allowedTransitions {
		if allowedStatus == to {
			return nil
		}
	}

	return NewInvalidStatusError("invalid status transition")
}

// Custom error types
type InvalidStatusError struct {
	message string
}

func NewInvalidStatusError(message string) *InvalidStatusError {
	return &InvalidStatusError{message: message}
}

func (e *InvalidStatusError) Error() string {
	return e.message
}
