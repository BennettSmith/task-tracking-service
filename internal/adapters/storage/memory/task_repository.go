package memory

import (
	"context"
	"fmt"
	"sync"
	"task-tracking-service/internal/core/domain"
	"task-tracking-service/pkg/errors"

	"github.com/google/uuid"
)

type TaskRepository struct {
	tasks map[string]*domain.Task
	mutex sync.RWMutex
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task.ID = uuid.New().String()
	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.NewNotFoundError(fmt.Sprintf("task with ID %s not found", id))
	}

	taskCopy := *task
	return &taskCopy, nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskCopy := *task
		tasks = append(tasks, &taskCopy)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("task with ID %s not found", task.ID))
	}

	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.NewNotFoundError(fmt.Sprintf("task with ID %s not found", id))
	}

	delete(r.tasks, id)
	return nil
}
