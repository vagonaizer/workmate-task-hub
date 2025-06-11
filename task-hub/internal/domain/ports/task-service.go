package ports

import (
	"time"

	"github.com/google/uuid"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
)

type TaskService interface {
	CreateTask(title, description string, priority models.TaskPriority, deadline time.Time) (*models.Task, error)
	GetTask(id uuid.UUID) (*models.Task, error)
	DeleteTask(id uuid.UUID) error
	ListTasks() ([]*models.Task, error)

	StartTask(id uuid.UUID) error
	CompleteTask(id uuid.UUID) error
	CancelTask(id uuid.UUID) error
	FailTask(id uuid.UUID) error

	UpdateTitle(id uuid.UUID, title string) error
	UpdateDescription(id uuid.UUID, description string) error
	UpdatePriority(id uuid.UUID, priority models.TaskPriority) error
	UpdateDeadline(id uuid.UUID, deadline time.Time) error
}
