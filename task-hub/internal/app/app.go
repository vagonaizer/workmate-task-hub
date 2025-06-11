package app

import (
	"github.com/gin-gonic/gin"
	"github.com/vagonaizer/workmate/task-hub/internal/config"
	"github.com/vagonaizer/workmate/task-hub/internal/domain/ports"
	inmemory "github.com/vagonaizer/workmate/task-hub/internal/repository/in-memory"
	service "github.com/vagonaizer/workmate/task-hub/internal/services/task-service"
	"github.com/vagonaizer/workmate/task-hub/internal/transport/http"
	"github.com/vagonaizer/workmate/task-hub/pkg/logger"
)

// App — основная структура приложения.
type App struct {
	Engine      *gin.Engine
	TaskService ports.TaskService
	Logger      *logger.Logger
}

// NewApp — собирает все зависимости и возвращает готовое приложение.
func NewApp(cfg *config.AppConfig) *App {
	// 1. Логгер
	logg := logger.NewLogger()

	// 2. Репозиторий (можно расширить switch для других хранилищ)
	var taskRepo ports.TaskRepository
	switch cfg.DB.Type {
	case config.DBInMemory:
		taskRepo = inmemory.NewInMemoryTaskRepository()
		logg.Info("Используется in-memory репозиторий")
	case config.DBPostgres:
		logg.Error("УПС, постгрес пока не поддерживается")
	default:
		logg.Error("Неизвестный тип репозитория: %v", cfg.DB.Type)
		panic("unknown repository type")
	}

	// 3. Сервис
	taskService := service.NewTaskService(taskRepo)

	// 4. Handler
	handler := http.NewHandler(taskService, logg)

	// 5. Gin + роуты
	engine := http.SetupRouter(handler)

	return &App{
		Engine:      engine,
		TaskService: taskService,
		Logger:      logg,
	}
}
