package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
)

type CreateTaskRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	Priority    models.TaskPriority `json:"priority"`
	Deadline    *time.Time          `json:"deadline,omitempty"`
}

type TaskResponse struct {
	ID          uuid.UUID           `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      models.TaskStatus   `json:"status"`
	Priority    models.TaskPriority `json:"priority"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	Duration    int64               `json:"duration_seconds"`
	Deadline    *time.Time          `json:"deadline,omitempty"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
