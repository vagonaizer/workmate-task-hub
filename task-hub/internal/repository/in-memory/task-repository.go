package inmemory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vagonaizer/workmate/task-hub/internal/common/apperror"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/ports"
)

// Убеждаемся, что InMemoryTaskRepository реализует интерфейс TaskRepository.
var _ ports.TaskRepository = (*InMemoryTaskRepository)(nil)

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*models.Task
}

// Конструктор.
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[uuid.UUID]*models.Task),
	}
}

func (r *InMemoryTaskRepository) Save(task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID()] = task
	return nil
}

func (r *InMemoryTaskRepository) GetByID(id uuid.UUID) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return nil, apperror.ErrRepoNotFound
	}
	return task, nil
}

func (r *InMemoryTaskRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return apperror.ErrRepoNotFound
	}
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryTaskRepository) List() ([]*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*models.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result, nil
}
