package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	task, err := NewTask("Test", "desc", TaskPriorityMedium)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Title() != "Test" {
		t.Errorf("expected title 'Test', got '%s'", task.Title())
	}
	if task.Status() != TaskStatusPending {
		t.Errorf("expected status 'pending', got '%s'", task.Status())
	}
}

func TestNewTask_EmptyTitle(t *testing.T) {
	_, err := NewTask("", "desc", TaskPriorityLow)
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestStartAndComplete(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	err := task.Start()
	assert.NoError(t, err)
	if task.Status() != TaskStatusInProgress {
		t.Errorf("expected status 'in_progress', got '%s'", task.Status())
	}
	err = task.Complete()
	assert.NoError(t, err)
	if task.Status() != TaskStatusCompleted {
		t.Errorf("expected status 'completed', got '%s'", task.Status())
	}
}

func TestCancel(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	if err := task.Cancel(); err != nil {
		t.Errorf("unexpected error on Cancel: %v", err)
	}
	if task.Status() != TaskStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", task.Status())
	}
}

func TestFail(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	err := task.Start()
	assert.NoError(t, err)
	err = task.Fail()
	assert.NoError(t, err)
	if task.Status() != TaskStatusFailed {
		t.Errorf("expected status 'failed', got '%s'", task.Status())
	}
}

func TestSetDeadline(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	future := time.Now().Add(24 * time.Hour)
	if err := task.SetDeadline(future); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !task.Deadline().Equal(future) {
		t.Errorf("deadline not set correctly")
	}
}

func TestSetDeadline_Past(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	past := time.Now().Add(-24 * time.Hour)
	if err := task.SetDeadline(past); err == nil {
		t.Error("expected error for past deadline")
	}
}

func TestSetTitleAndDescription(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	if err := task.SetTitle("NewTitle"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.Title() != "NewTitle" {
		t.Errorf("title not updated")
	}
	task.SetDescription("new desc")
	if task.Description() != "new desc" {
		t.Errorf("description not updated")
	}
}

func TestSetPriority(t *testing.T) {
	task, _ := NewTask("Test", "desc", TaskPriorityLow)
	if err := task.SetPriority(TaskPriorityHigh); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.Priority() != TaskPriorityHigh {
		t.Errorf("priority not updated")
	}
	if err := task.SetPriority("invalid"); err == nil {
		t.Error("expected error for invalid priority")
	}
}
