package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"task-tracking-service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Parse command line flags
	var (
		command       = flag.String("command", "", "migrate command (up/down)")
		migrationsDir = flag.String("migrations", "migrations", "migrations directory")
	)
	flag.Parse()

	if *command == "" {
		log.Fatal("command is required")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Construct database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Get absolute path to migrations
	absPath, err := filepath.Abs(*migrationsDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create migrate instance
	m, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Execute migration command
	switch *command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Successfully ran migrations")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Successfully rolled back migrations")

	default:
		log.Fatalf("Invalid command: %s", *command)
	}
}
