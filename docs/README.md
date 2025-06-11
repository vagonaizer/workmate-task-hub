## Domain Layer (models)

В доменной области определена главная структура - Task, представляющая из себя задачу в нашей программе.

```
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

```
Следуя принципам чистой архитектуры важно избегать прямого вмешательства в поля структуры из других слоев, именно поэтому поля в этой структуре -
- неэкспортируемые. Поля структуры изменяются за счет методов-сеттеров. 

Пример:

Изменения статуса задачи в случае невыполнения задачи на "failed".
```
func (t *Task) Fail() error {
	if t.status != TaskStatusInProgress {
		return fmt.Errorf("%w: %v", ErrInvalidStatus, t.status)
	}
	t.status = TaskStatusFailed
	t.updatedAt = time.Now()
	return nil
}

```

Для этих методов написаны юниты в той же директории где расположена entity сущность.

---

## Domain layer (ports)

Идея создания подобной директории и использование интерфейсов не по месту назначения возникла после столкновения с Hexagonal Architecture.
Такой подход мне нравится больше, в особенности для построения большого монолита.

Интерфейс репозитория:

```

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

```

Интерфейс сервиса (usecase):

```

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

```

## Repository

Слой репозитория отвечает за простейшую обработку данных: добавление/удаление/отображение всех записей или по ID 
без затрагивания сложных процессов валидации или переброса ненужного функционала место которому на сервисном уровне. Иначе мы рискуем столкнуться с аномалиями, что не есть хорошо.

```

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*models.Task
}

```

Репозиторий in-memory. Потокобезопасной за счет использования мьютекса. 
Испольщуется структура данных [id задачи] -> структура task. 

Так же написаны unit тесты.

## Service (usecase) layer:

TODO: дописать Readme
 
