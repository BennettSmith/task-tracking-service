package factory

import (
	"testing"

	"task-tracking-service/internal/adapters/storage/memory"
	"task-tracking-service/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryFactory_CreateTaskRepository(t *testing.T) {
	tests := []struct {
		name           string
		repositoryType string
		expectError    bool
	}{
		{
			name:           "create memory repository",
			repositoryType: string(MemoryRepository),
			expectError:    false,
		},
		{
			name:           "create postgres repository",
			repositoryType: string(PostgresRepository),
			expectError:    true,
		},
		{
			name:           "unknown repository type",
			repositoryType: "unknown",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config
			cfg := &config.Config{
				Repository: config.RepositoryConfig{
					Type: tt.repositoryType,
				},
			}

			// Create factory
			factory := NewRepositoryFactory(cfg)

			// Create repository
			repo, err := factory.CreateTaskRepository()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, repo)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, repo)

			// Verify repository type
			switch tt.repositoryType {
			case string(MemoryRepository):
				assert.IsType(t, &memory.TaskRepository{}, repo)
			}
		})
	}
}

func TestNewRepositoryFactory(t *testing.T) {
	cfg := &config.Config{}
	factory := NewRepositoryFactory(cfg)

	assert.NotNil(t, factory)
	assert.Equal(t, cfg, factory.config)
}
