package http

import (
	"net/http"
	"task-tracking-service/internal/core/domain"
	"task-tracking-service/internal/core/services"
	"time"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	task, err := h.taskService.CreateTask(c.Request().Context(), req.Title, req.Description, req.DueDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create task")
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c echo.Context) error {
	id := c.Param("id")
	task, err := h.taskService.GetTask(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ListTasks(c echo.Context) error {
	tasks, err := h.taskService.ListTasks(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch tasks")
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	id := c.Param("id")
	var task domain.Task
	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	task.ID = id
	updatedTask, err := h.taskService.UpdateTask(c.Request().Context(), &task)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update task")
	}

	return c.JSON(http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id := c.Param("id")
	if err := h.taskService.DeleteTask(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete task")
	}

	return c.NoContent(http.StatusNoContent)
}
