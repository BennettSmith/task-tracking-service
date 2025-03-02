package factory

import (
	"fmt"

	"task-tracking-service/internal/adapters/storage/memory"
	"task-tracking-service/internal/config"
	"task-tracking-service/internal/core/ports"
)

// RepositoryType defines the available repository implementations
type RepositoryType string

const (
	MemoryRepository   RepositoryType = "memory"
	PostgresRepository RepositoryType = "postgres"
)

// RepositoryFactory creates and configures repositories
type RepositoryFactory struct {
	cfg *config.Config
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(cfg *config.Config) *RepositoryFactory {
	return &RepositoryFactory{
		cfg: cfg,
	}
}

// CreateTaskRepository creates a task repository based on configuration
func (f *RepositoryFactory) CreateTaskRepository() (ports.TaskRepository, error) {
	repoType := RepositoryType(f.cfg.Repository.Type)

	switch repoType {
	case MemoryRepository:
		return memory.NewTaskRepository(), nil

	case PostgresRepository:
		return nil, fmt.Errorf("postgres repository not implemented yet")

	default:
		return nil, fmt.Errorf("unknown repository type: %s", repoType)
	}
}
