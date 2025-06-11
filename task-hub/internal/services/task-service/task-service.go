package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/vagonaizer/workmate/task-hub/internal/common/apperror"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/ports"
)

// TaskService — сервисный слой для работы с задачами.

// Убеждаемся, что TaskService реализует интерфейс TaskService.
var _ ports.TaskService = (*TaskService)(nil)

type TaskService struct {
	repo ports.TaskRepository
}

// Конструктор принимающий на вход репозиторий.
func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask — создание новой задачи.
func (s *TaskService) CreateTask(title, description string, priority models.TaskPriority, deadline time.Time) (*models.Task, error) {
	task, err := models.NewTask(title, description, priority)
	if err != nil {
		return nil, apperror.ErrServiceValidation
	}
	if !deadline.IsZero() {
		if err := task.SetDeadline(deadline); err != nil {
			return nil, apperror.ErrServiceValidation
		}
	}
	if err := s.repo.Save(task); err != nil {
		return nil, apperror.ErrRepoSaveFailed
	}
	return task, nil
}

// GetTask — получение задачи по идентификатору.
func (s *TaskService) GetTask(id uuid.UUID) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.ErrRepoNotFound
	}
	return task, nil
}

// DeleteTask — удаление задачи по идентификатору.
func (s *TaskService) DeleteTask(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// ListTasks — получение списка всех задач.
func (s *TaskService) ListTasks() ([]*models.Task, error) {
	return s.repo.List()
}

// StartTask — переводит задачу в статус "в работе".
func (s *TaskService) StartTask(id uuid.UUID) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.Start(); err != nil {
		return err
	}
	return s.repo.Save(task)
}

// CompleteTask — завершает задачу.
func (s *TaskService) CompleteTask(id uuid.UUID) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.Complete(); err != nil {
		return err
	}
	return s.repo.Save(task)
}

// CancelTask — отменяет задачу.
func (s *TaskService) CancelTask(id uuid.UUID) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.Cancel(); err != nil {
		return err
	}
	return s.repo.Save(task)
}

// FailTask — переводит задачу в статус "ошибка при выполнении".
func (s *TaskService) FailTask(id uuid.UUID) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.Fail(); err != nil {
		return err
	}
	return s.repo.Save(task)
}

// DeleteDomainTask — переводит задачу в статус "удалена" (soft delete).
func (s *TaskService) DeleteDomainTask(id uuid.UUID) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.Delete(); err != nil {
		return err
	}
	return s.repo.Save(task)
}

// SetDeadline — устанавливает дедлайн задачи.
func (s *TaskService) SetDeadline(id uuid.UUID, deadline time.Time) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.SetDeadline(deadline); err != nil {
		return err
	}
	return s.repo.Save(task)
}

// UpdateTitle — изменяет заголовок задачи.
func (s *TaskService) UpdateTitle(id uuid.UUID, title string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.SetTitle(title); err != nil {
		return apperror.ErrServiceValidation
	}
	return s.repo.Save(task)
}

// UpdateDescription — изменяет описание задачи.
func (s *TaskService) UpdateDescription(id uuid.UUID, description string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	task.SetDescription(description)
	return s.repo.Save(task)
}

// UpdatePriority — изменяет приоритет задачи.
func (s *TaskService) UpdatePriority(id uuid.UUID, priority models.TaskPriority) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	// Добавь метод SetPriority в models.Task, если его нет
	if err := task.SetPriority(priority); err != nil {
		return apperror.ErrServiceValidation
	}
	return s.repo.Save(task)
}

// UpdateDeadline — изменяет дедлайн задачи.
func (s *TaskService) UpdateDeadline(id uuid.UUID, deadline time.Time) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.ErrRepoNotFound
	}
	if err := task.SetDeadline(deadline); err != nil {
		return apperror.ErrServiceValidation
	}
	return s.repo.Save(task)
}
