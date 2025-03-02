package factory

import (
	"database/sql"
	"fmt"

	"task-tracking-service/internal/adapters/storage/memory"
	"task-tracking-service/internal/adapters/storage/postgres"
	"task-tracking-service/internal/config"
	"task-tracking-service/internal/core/ports"
	"task-tracking-service/internal/migrations"
)

// RepositoryType defines the available repository implementations
type RepositoryType string

const (
	MemoryRepository   RepositoryType = "memory"
	PostgresRepository RepositoryType = "postgres"
)

// RepositoryFactory creates and configures repositories
type RepositoryFactory struct {
	config *config.Config
	db     *sql.DB
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(config *config.Config) *RepositoryFactory {
	return &RepositoryFactory{
		config: config,
	}
}

// CreateTaskRepository creates a task repository based on configuration
func (f *RepositoryFactory) CreateTaskRepository() (ports.TaskRepository, error) {
	switch f.config.Repository.Type {
	case "postgres":
		if f.db == nil {
			// Initialize database connection
			db, err := postgres.NewDB(f.config.Database)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize database: %w", err)
			}
			f.db = db

			// Run migrations
			if err := migrations.MigrateDB(db, "migrations"); err != nil {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		}
		return postgres.NewTaskRepository(f.db), nil

	case "memory":
		return memory.NewTaskRepository(), nil

	default:
		return nil, fmt.Errorf("unknown repository type: %s", f.config.Repository.Type)
	}
}

// Close cleans up any resources (like database connections)
func (f *RepositoryFactory) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}
