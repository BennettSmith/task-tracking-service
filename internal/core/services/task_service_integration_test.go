package services

import (
	"context"
	"testing"
	"time"

	"task-tracking-service/internal/adapters/storage/memory"
	"task-tracking-service/internal/core/domain"

	"github.com/stretchr/testify/suite"
)

type TaskServiceIntegrationSuite struct {
	suite.Suite
	ctx     context.Context
	service *TaskService
	repo    *memory.TaskRepository
}

func (s *TaskServiceIntegrationSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = memory.NewTaskRepository()
	s.service = NewTaskService(s.repo)
}

func TestTaskServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceIntegrationSuite))
}

func (s *TaskServiceIntegrationSuite) TestTaskLifecycle() {
	// Create a task
	title := "Integration Test Task"
	desc := "Testing full task lifecycle"
	dueDate := time.Now().Add(24 * time.Hour)

	task, err := s.service.CreateTask(s.ctx, title, desc, dueDate)
	s.NoError(err)
	s.NotEmpty(task.ID)
	s.Equal(title, task.Title)
	s.Equal(desc, task.Description)
	s.Equal(domain.StatusPending, task.Status)

	// Get the task
	retrieved, err := s.service.GetTask(s.ctx, task.ID)
	s.NoError(err)
	s.Equal(task.ID, retrieved.ID)
	s.Equal(task.Title, retrieved.Title)

	// Update the task
	retrieved.Status = domain.StatusInProgress
	retrieved.Title = "Updated Title"
	updated, err := s.service.UpdateTask(s.ctx, retrieved)
	s.NoError(err)
	s.Equal("Updated Title", updated.Title)
	s.Equal(domain.StatusInProgress, updated.Status)

	// List tasks
	tasks, err := s.service.ListTasks(s.ctx)
	s.NoError(err)
	s.Len(tasks, 1)
	s.Equal(updated.ID, tasks[0].ID)

	// Delete the task
	err = s.service.DeleteTask(s.ctx, task.ID)
	s.NoError(err)

	// Verify deletion
	tasks, err = s.service.ListTasks(s.ctx)
	s.NoError(err)
	s.Empty(tasks)
}

func (s *TaskServiceIntegrationSuite) TestInvalidStatusTransitions() {
	// Create a task
	task, err := s.service.CreateTask(s.ctx, "Test Task", "Description", time.Now().Add(24*time.Hour))
	s.NoError(err)

	// Try invalid status transition
	task.Status = "invalid-status"
	_, err = s.service.UpdateTask(s.ctx, task)
	s.Error(err)
	s.IsType(&InvalidStatusError{}, err)
}

func (s *TaskServiceIntegrationSuite) TestConcurrentOperations() {
	// Create initial task
	task, err := s.service.CreateTask(s.ctx, "Concurrent Test", "Description", time.Now().Add(24*time.Hour))
	s.NoError(err)

	// Simulate concurrent updates
	done := make(chan bool)
	go func() {
		task.Title = "Update 1"
		_, err := s.service.UpdateTask(s.ctx, task)
		s.NoError(err)
		done <- true
	}()

	go func() {
		task.Title = "Update 2"
		_, err := s.service.UpdateTask(s.ctx, task)
		s.NoError(err)
		done <- true
	}()

	// Wait for both updates
	<-done
	<-done

	// Verify final state
	updated, err := s.service.GetTask(s.ctx, task.ID)
	s.NoError(err)
	s.NotEqual("Concurrent Test", updated.Title)
}
