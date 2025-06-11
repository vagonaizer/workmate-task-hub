package http

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/ports"
	"github.com/vagonaizer/workmate/task-hub/pkg/logger"
)

type Handler struct {
	taskService ports.TaskService
	logger      *logger.Logger
}

func NewHandler(taskService ports.TaskService, logger *logger.Logger) *Handler {
	return &Handler{taskService: taskService, logger: logger}
}

// @@route POST /api/tasks
// @@desc  Создать задачу
// @@accept json
// @@success 201 TaskResponse
// @@error 400 Ошибка валидации или данных
func (h *Handler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var deadline time.Time
	if req.Deadline != nil {
		deadline = *req.Deadline
	}
	task, err := h.taskService.CreateTask(req.Title, req.Description, req.Priority, deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем задачу в файл
	taskResp := toTaskResponse(task)
	if err := saveTaskToFile(taskResp); err != nil {
		h.logger.Error("Ошибка сохранения задачи в файл: %v", err)
	}

	c.JSON(http.StatusCreated, taskResp)
}

// @@route GET /api/tasks/:id
// @@desc  Получить задачу по id
// @@success 200 TaskResponse
// @@error 400 invalid id
// @@error 404 not found
func (h *Handler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := h.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toTaskResponse(task))
}

// @@route DELETE /api/tasks/:id
// @@desc  Удалить задачу по id
// @@success 204
// @@error 400 invalid id
// @@error 404 not found
func (h *Handler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.taskService.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @@route GET /api/tasks
// @@desc  Получить список всех задач
// @@success 200 TaskListResponse
func (h *Handler) ListTasks(c *gin.Context) {
	tasks, err := h.taskService.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp = append(resp, toTaskResponse(t))
	}
	c.JSON(http.StatusOK, TaskListResponse{Tasks: resp})
}

// @@route PATCH /api/tasks/:id/status
// @@desc  Изменить статус задачи
// @@accept json
// @@success 204
// @@error 400 invalid id
// @@error 400 invalid status
// @@error 404 not found
func (h *Handler) UpdateTaskStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Status models.TaskStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Изменение статуса задачи с id: " + id.String() + " на " + string(req.Status))
	switch req.Status {
	case models.TaskStatusInProgress:
		err = h.taskService.StartTask(id)
	case models.TaskStatusCompleted:
		err = h.taskService.CompleteTask(id)
	case models.TaskStatusCancelled:
		err = h.taskService.CancelTask(id)
	case models.TaskStatusFailed:
		err = h.taskService.FailTask(id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "id": idStr})
		return
	}
	c.Status(http.StatusNoContent)
}

// @@route GET /api/tasks/:id/status
// @@desc  Получить статус задачи
// @@success 200 {"status": string}
// @@error 400 invalid id
// @@error 404 not found
func (h *Handler) GetTaskStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	h.logger.Info("Получение статуса задачи с id: " + id.String())
	task, err := h.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "id": idStr})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": task.Status()})
}

// @@route PATCH /api/tasks/:id/title
// @@desc  Изменить название задачи
// @@accept json
// @@success 204
// @@error 400 invalid id
// @@error 400 Ошибка валидации
// @@error 404 not found
func (h *Handler) UpdateTaskTitle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Изменение названия задачи с id: %s на %s", id.String(), req.Title)
	if err := h.taskService.UpdateTitle(id, req.Title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "id": idStr})
		return
	}
	c.Status(http.StatusNoContent)
}

// @@route PATCH /api/tasks/:id/description
// @@desc  Изменить описание задачи
// @@accept json
// @@success 204
// @@error 400 invalid id
// @@error 400 Ошибка валидации
// @@error 404 not found
func (h *Handler) UpdateTaskDescription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Изменение описания задачи с id: " + id.String())
	if err := h.taskService.UpdateDescription(id, req.Description); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "id": idStr})
		return
	}
	c.Status(http.StatusNoContent)
}

// toTaskResponse — маппинг доменной задачи в DTO
func toTaskResponse(t *models.Task) TaskResponse {
	var completedAt *time.Time
	if !t.CompletedAt().IsZero() {
		completedAt = &[]time.Time{t.CompletedAt()}[0]
	}
	var deadline *time.Time
	if !t.Deadline().IsZero() {
		deadline = &[]time.Time{t.Deadline()}[0]
	}
	return TaskResponse{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Status:      t.Status(),
		Priority:    t.Priority(),
		CreatedAt:   t.CreatedAt(),
		UpdatedAt:   t.UpdatedAt(),
		CompletedAt: completedAt,
		Duration:    int64(t.Duration().Seconds()),
		Deadline:    deadline,
	}
}

// saveTaskToFile сохраняет задачу в examples/tasks.json
func saveTaskToFile(task TaskResponse) error {
	const filePath = "examples/tasks.json"
	var tasks []TaskResponse

	// Прочитать существующий файл, если есть
	if data, err := os.ReadFile(filePath); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &tasks); err != nil {
			return err
		}
	}

	// Добавить новую задачу
	tasks = append(tasks, task)

	// Сохранить обратно
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
