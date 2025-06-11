package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/api")
	// Маршруты для работы с задачами
	tasks := api.Group("/tasks")
	{
		tasks.POST("", handler.CreateTask)
		tasks.GET("", handler.ListTasks)
		tasks.GET("/:id", handler.GetTask)
		tasks.DELETE("/:id", handler.DeleteTask)
		tasks.PATCH("/:id/status", handler.UpdateTaskStatus)
		tasks.GET("/:id/status", handler.GetTaskStatus)
		tasks.PATCH("/:id/title", handler.UpdateTaskTitle)
		tasks.PATCH("/:id/description", handler.UpdateTaskDescription)
	}
	return router
}
