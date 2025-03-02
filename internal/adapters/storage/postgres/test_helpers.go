package postgres

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"

	"your-module/internal/adapters/storage/postgres/migrations"
	"your-module/internal/config"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Get test database configuration
	dbConfig := getTestDBConfig()

	// Connect to database
	db, err := NewDB(dbConfig)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Get path to migrations
	_, b, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(b), "../../../../migrations")

	// Run migrations
	if err := migrations.MigrateDB(db, migrationsPath); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Clean up database
	if _, err := db.Exec("TRUNCATE TABLE tasks"); err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	return db
}

func getTestDBConfig() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Name:     "taskdb_test",
		SSLMode:  "disable",
	}
}
