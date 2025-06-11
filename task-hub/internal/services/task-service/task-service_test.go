package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
	"github.com/vagonaizer/workmate/task-hub/internal/services/task-service/mocks"
)

func TestTaskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := NewTaskService(mockRepo)

	title := "Test"
	desc := "desc"
	priority := models.TaskPriorityHigh
	deadline := time.Now().Add(24 * time.Hour)

	task, _ := models.NewTask(title, desc, priority)
	err := task.SetDeadline(deadline)
	assert.NoError(t, err)
	mockRepo.EXPECT().Save(gomock.Any()).Return(nil)

	created, err := service.CreateTask(title, desc, priority, deadline)
	assert.NoError(t, err)
	assert.Equal(t, title, created.Title())
	assert.Equal(t, desc, created.Description())
	assert.Equal(t, priority, created.Priority())
}

func TestTaskService_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := NewTaskService(mockRepo)

	id := uuid.New()
	task, _ := models.NewTask("Test", "desc", models.TaskPriorityLow)
	mockRepo.EXPECT().GetByID(id).Return(task, nil)

	got, err := service.GetTask(id)
	assert.NoError(t, err)
	assert.Equal(t, task, got)
}

func TestTaskService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := NewTaskService(mockRepo)

	id := uuid.New()
	mockRepo.EXPECT().Delete(id).Return(nil)

	err := service.DeleteTask(id)
	assert.NoError(t, err)
}

func TestTaskService_ListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := NewTaskService(mockRepo)
	task1, _ := models.NewTask("T1", "desc1", models.TaskPriorityLow)
	task2, _ := models.NewTask("T2", "desc2", models.TaskPriorityHigh)
	mockRepo.EXPECT().List().Return([]*models.Task{task1, task2}, nil)

	list, err := service.ListTasks()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := NewTaskService(mockRepo)
	id := uuid.New()
	mockRepo.EXPECT().GetByID(id).Return(nil, assert.AnError)

	_, err := service.GetTask(id)
	assert.Error(t, err)
}
