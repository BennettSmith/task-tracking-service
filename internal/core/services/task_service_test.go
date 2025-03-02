package services

import (
	"context"
	"testing"
	"time"

	"task-tracking-service/internal/core/domain"
	"task-tracking-service/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository is a mock implementation of ports.TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	if len(args) > 0 && args.Get(0) != nil {
		task.ID = "test-id"
	}
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successfully creates task", func(t *testing.T) {
		title := "Test Task"
		description := "Test Description"
		dueDate := time.Now().Add(24 * time.Hour)

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Task")).Return(nil)

		task, err := service.CreateTask(ctx, title, description, dueDate)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, title, task.Title)
		assert.Equal(t, description, task.Description)
		assert.Equal(t, domain.StatusPending, task.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskService_GetTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successfully gets existing task", func(t *testing.T) {
		expectedTask := &domain.Task{
			ID:          "test-id",
			Title:       "Test Task",
			Description: "Test Description",
			Status:      domain.StatusPending,
		}

		mockRepo.On("GetByID", ctx, "test-id").Return(expectedTask, nil)

		task, err := service.GetTask(ctx, "test-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error for non-existent task", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, errors.NewNotFoundError("task not found"))

		task, err := service.GetTask(ctx, "non-existent")

		assert.Error(t, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successfully updates task", func(t *testing.T) {
		existingTask := &domain.Task{
			ID:        "test-id",
			Title:     "Test Task",
			Status:    domain.StatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		updatedTask := &domain.Task{
			ID:        "test-id",
			Title:     "Updated Test Task",
			Status:    domain.StatusInProgress,
			CreatedAt: existingTask.CreatedAt,
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetByID", ctx, "test-id").Return(existingTask, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Task")).Return(nil)

		result, err := service.UpdateTask(ctx, updatedTask)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedTask.Title, result.Title)
		assert.Equal(t, domain.StatusInProgress, result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fails with invalid status transition", func(t *testing.T) {
		existingTask := &domain.Task{
			ID:     "test-id",
			Status: domain.StatusPending,
		}

		invalidTask := &domain.Task{
			ID:     "test-id",
			Status: "invalid-status",
		}

		mockRepo.On("GetByID", ctx, "test-id").Return(existingTask, nil)

		result, err := service.UpdateTask(ctx, invalidTask)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.IsType(t, &InvalidStatusError{}, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("successfully deletes existing task", func(t *testing.T) {
		existingTask := &domain.Task{
			ID: "test-id",
		}

		mockRepo.On("GetByID", ctx, "test-id").Return(existingTask, nil)
		mockRepo.On("Delete", ctx, "test-id").Return(nil)

		err := service.DeleteTask(ctx, "test-id")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fails to delete non-existent task", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, errors.NewNotFoundError("task not found"))

		err := service.DeleteTask(ctx, "non-existent")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
