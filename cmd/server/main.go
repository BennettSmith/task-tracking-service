package main

import (
	"log"
	"task-tracking-service/internal/adapters/http"
	"task-tracking-service/internal/adapters/storage/factory"
	"task-tracking-service/internal/config"
	"task-tracking-service/internal/core/services"
	// You'll need to import your repository implementation once it's created
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create repository factory
	repoFactory := factory.NewRepositoryFactory(cfg)
	defer repoFactory.Close()

	// Create task repository using the factory
	taskRepo, err := repoFactory.CreateTaskRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Initialize service with the repository from factory
	taskService := services.NewTaskService(taskRepo)

	// Initialize handlers
	taskHandler := http.NewTaskHandler(taskService)

	// Setup router
	router := http.NewRouter(taskHandler)

	// Start server
	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Start(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
