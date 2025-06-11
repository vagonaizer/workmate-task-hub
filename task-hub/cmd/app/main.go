package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vagonaizer/workmate/task-hub/internal/app"
	"github.com/vagonaizer/workmate/task-hub/internal/config"
	"github.com/vagonaizer/workmate/task-hub/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	application := app.NewApp(cfg)
	application.Logger.Info("Запуск %s v%s на :%s", cfg.AppName, cfg.AppVersion, cfg.Server.Port)
	hireMe()

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: application.Engine,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			application.Logger.Error("Ошибка запуска сервера: %v", err)
		}
	}()

	// Graceful shutdown
	waitForShutdown(server, application.Logger)
	os.Exit(0)
}

func waitForShutdown(server *http.Server, logger *logger.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Получен сигнал завершения: %v, останавливаем сервер...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Ошибка graceful shutdown: %v", err)
	} else {
		logger.Info("Graceful shutdown завершен ( =ω= )")
	}
}

func hireMe() {
	fmt.Println(`  _   _ ___ ____  _____ 
 | | | |_ _|  _ \| ____|
 | |_| || || |_) |  _|  
 |  _  || ||  _ <| |___ 
 |_| |_|___|_|_\_\_____|
      |  \/  | ____|    
      | |\/| |  _|      
      | |  | | |___     
    __|_| _|_|_____|    
   |  _ \| |   / ___|   
   | |_) | |   \___ \   
   |  __/| |___ ___) |  
   |_|   |_____|____/   
                       `)
}
