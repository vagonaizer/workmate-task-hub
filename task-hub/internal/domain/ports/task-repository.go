package ports

import (
	"github.com/google/uuid"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/models"
)

/*
-- Интерфейс-контракт для репозитория

	Сначала хотел разместить интерфейс в месте его реализации
	на сервисном (usecase) уровне, но мне ближе подход из
	Hexagonal Architecture с отдельным объявлением интерфейсов
	в доменной области. Он мне кажется лаконичнее, в особенности
	учитывая реалии больших тяжеловесных проектов, если монолит так вообще
	отдельная директория под интерфейсы - отличная практика.
*/

type TaskRepository interface {
	// Save сохраняет задачу (создаёт новую или обновляет существующую).
	Save(task *models.Task) error

	// GetByID возвращает задачу по её идентификатору.
	GetByID(id uuid.UUID) (*models.Task, error)

	// Delete удаляет задачу по идентификатору.
	Delete(id uuid.UUID) error

	// List возвращает все задачи.
	List() ([]*models.Task, error)
}
