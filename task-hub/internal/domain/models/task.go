package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Алиасы для статусов и приоритов задач
type (
	TaskStatus   string
	TaskPriority string
)

// Можно было бы использовать iota, в реальном проекте
// использовал бы именно перечисления из-за скорости выполнения,
// но в данном случае, для демонстрации использую строковые константы.
const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusDeleted    TaskStatus = "deleted"
)

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

// Фундаментальная сущность Task -- представляет собой задачу.
// -- 1. Мы используем неэкспортирумые поля, доступ к ним осуществляется через методы core-logic.
// -- 2. Мы не используем json-теги, для маршалинга используются ДТО.
type Task struct {
	id          uuid.UUID     // ID задачи
	title       string        // Заголовок-название задачи
	description string        // Описание задачи, может быть пустым
	status      TaskStatus    // Статус задачи
	priority    TaskPriority  // Приоритет задачи
	createdAt   time.Time     // Время создания задачи
	updatedAt   time.Time     // Время последнего обновления задачи: статус, описание и т.д.
	completedAt time.Time     // Время заверешения задачи
	duration    time.Duration // Сколько ушло времени на выполнение задачи
	deadline    time.Time     // Крайний срок выполнения задачи
}

// Конструктор Task.
func NewTask(title, description string, priority TaskPriority) (*Task, error) {
	// Логика следующая:
	// 1. Задача без title (названия) -- не может существовать.
	// 2. Названия задачи способно полноценно описать задачу -- нет необходимости в description.
	if title == "" {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTitle, title)
	}
	// 3. Приоритет задачи (если не выставлен явно, по умолчанию -- low)
	// Не будем рушить всю бизнес-логику, если приоритет не выставлен.
	if priority == "" {
		priority = TaskPriorityLow
	}
	return &Task{
		id:          uuid.New(),
		title:       title,
		description: description,
		status:      TaskStatusPending,
		priority:    priority,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		completedAt: time.Time{},
		duration:    0,
		deadline:    time.Time{},
	}, nil
}

// Сеттеры - методы для изменения состояния task.
// -- 1. Старт задачи
// -- 2. Завершение задачи
// -- 3. Отмена задачи
// -- 4. Смена статуса на failed
// -- 5. Удаление задачи
// -- 6. Установка дедлайна
// -- 7. Изменение заголовка (названия) задачи
// -- 8. Задает описание задачи
// -- 9. Изменение приоритета задачи

// Start — переводит задачу в статус "в работе".
// -- 1. Проверяет, что задача в статусе "pending".
// -- 2. Устанавливает статус "in_progress" и обновляет время обновления.
func (t *Task) Start() error {
	if t.status != TaskStatusPending {
		return ErrInvalidStatus
	}
	t.status = TaskStatusInProgress
	t.updatedAt = time.Now()
	return nil
}

// Complete — завершает задачу.
// -- 1. Проверяет, что задача в статусе "in_progress".
// -- 2. Устанавливает статус "completed", фиксирует время завершения и считает duration.
func (t *Task) Complete() error {
	if t.status != TaskStatusInProgress {
		return fmt.Errorf("%w: %v", ErrInvalidStatus, t.status)
	}
	t.status = TaskStatusCompleted
	t.completedAt = time.Now()
	t.updatedAt = t.completedAt
	t.duration = t.completedAt.Sub(t.createdAt)
	return nil
}

// Cancel — отменяет задачу.
// -- 1. Проверяет, что задача в статусе "pending" или "in_progress".
// -- 2. Устанавливает статус "cancelled" и обновляет время обновления.
func (t *Task) Cancel() error {
	if t.status != TaskStatusPending && t.status != TaskStatusInProgress {
		return fmt.Errorf("%w: %v", ErrInvalidStatus, t.status)
	}
	t.status = TaskStatusCancelled
	t.updatedAt = time.Now()
	return nil
}

// Fail — переводит задачу в статус "failed".
// -- 1. Проверяет, что задача в статусе "in_progress".
// -- 2. Устанавливает статус "failed" и обновляет время обновления.
func (t *Task) Fail() error {
	if t.status != TaskStatusInProgress {
		return fmt.Errorf("%w: %v", ErrInvalidStatus, t.status)
	}
	t.status = TaskStatusFailed
	t.updatedAt = time.Now()
	return nil
}

// Delete — переводит задачу в статус "deleted".
// -- 1. Устанавливает статус "deleted" и обновляет время обновления.
func (t *Task) Delete() error {
	t.status = TaskStatusDeleted
	t.updatedAt = time.Now()
	return nil
}

// SetDeadline — изменяет дедлайн задачи.
// -- 1. Проверяет, что дедлайн не в прошлом.
// -- 2. Не позволяет менять дедлайн для завершённых, отменённых и удалённых задач.
// -- 3. Обновляет поле и время обновления.
func (t *Task) SetDeadline(deadline time.Time) error {
	if !deadline.IsZero() && deadline.Before(time.Now()) {
		return fmt.Errorf("%w: deadline must be in the future", ErrInvalidDeadline)
	}
	if t.status == TaskStatusCompleted || t.status == TaskStatusCancelled || t.status == TaskStatusDeleted {
		return fmt.Errorf("%w: cannot set deadline for finished/cancelled/deleted task", ErrInvalidStatus)
	}
	t.deadline = deadline
	t.updatedAt = time.Now()
	return nil
}

// SetTitle — изменяет заголовок (название) задачи.
// -- 1. Проверяет, что заголовок не пустой.
// -- 2. Обновляет поле и время обновления.
func (t *Task) SetTitle(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}
	t.title = title
	t.updatedAt = time.Now()
	return nil
}

// SetDescription — изменяет описание задачи.
// -- 1. Обновляет поле и время обновления.
func (t *Task) SetDescription(description string) {
	t.description = description
	t.updatedAt = time.Now()
}

// SetPriority — изменяет приоритет задачи.
// -- 1. Проверяет валидность приоритета.
// -- 2. Обновляет поле и время обновления.
func (t *Task) SetPriority(priority TaskPriority) error {
	if priority != TaskPriorityLow && priority != TaskPriorityMedium && priority != TaskPriorityHigh {
		return fmt.Errorf("%w: %v", ErrInvalidPriority, priority)
	}
	t.priority = priority
	t.updatedAt = time.Now()
	return nil
}

// Геттеры

func (t *Task) ID() uuid.UUID {
	return t.id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func (t *Task) Priority() TaskPriority {
	return t.priority
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Task) CompletedAt() time.Time {
	return t.completedAt
}

func (t *Task) Duration() time.Duration {
	return t.duration
}

func (t *Task) Deadline() time.Time {
	return t.deadline
}
