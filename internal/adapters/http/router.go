package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(taskHandler *TaskHandler) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Task routes
	tasks := v1.Group("/task")
	tasks.POST("", taskHandler.CreateTask)
	tasks.GET("", taskHandler.ListTasks)
	tasks.GET("/:id", taskHandler.GetTask)
	tasks.PUT("/:id", taskHandler.UpdateTask)
	tasks.DELETE("/:id", taskHandler.DeleteTask)

	return e
}
