package inmemory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
)

func TestInMemoryTaskRepository_SaveAndGetByID(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	task, err := models.NewTask("Test", "desc", models.TaskPriorityLow)
	assert.NoError(t, err)

	err = repo.Save(task)
	assert.NoError(t, err)

	got, err := repo.GetByID(task.ID())
	assert.NoError(t, err)
	assert.Equal(t, task, got)
}

func TestInMemoryTaskRepository_GetByID_NotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	_, err := repo.GetByID(uuid.New())
	assert.Error(t, err)
}

func TestInMemoryTaskRepository_Delete(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	task, _ := models.NewTask("Test", "desc", models.TaskPriorityLow)
	err := repo.Save(task)
	assert.NoError(t, err)
	err = repo.Delete(task.ID())
	assert.NoError(t, err)

	_, err = repo.GetByID(task.ID())
	assert.Error(t, err)
}

func TestInMemoryTaskRepository_Delete_NotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	err := repo.Delete(uuid.New())
	assert.Error(t, err)
}

func TestInMemoryTaskRepository_List(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	task1, _ := models.NewTask("T1", "desc1", models.TaskPriorityLow)
	task2, _ := models.NewTask("T2", "desc2", models.TaskPriorityHigh)
	err := repo.Save(task1)
	assert.NoError(t, err)
	err = repo.Save(task2)
	assert.NoError(t, err)

	list, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Contains(t, list, task1)
	assert.Contains(t, list, task2)
}
