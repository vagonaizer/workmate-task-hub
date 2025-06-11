.PHONY: all build run clean test lint help

APP_NAME=task-hub
CMD_PATH=task-hub/cmd/app/main.go
BIN_PATH=bin/$(APP_NAME)

RESET=\033[0m
BOLD=\033[1m
GREEN=\033[32m
YELLOW=\033[33m
BLUE=\033[34m
CYAN=\033[36m

all: build

help:
	@echo "\n$(BOLD)$(CYAN)Доступные команды:$(RESET)"
	@echo "  $(BOLD)make build$(RESET)   — 🛠️  Собрать бинарник"
	@echo "  $(BOLD)make run$(RESET)     — 🚀  Запустить приложение"
	@echo "  $(BOLD)make test$(RESET)    — 🧪  Прогнать все тесты"
	@echo "  $(BOLD)make lint$(RESET)    — 🔍  Запустить линтер"
	@echo "  $(BOLD)make clean$(RESET)   — 🧹  Очистить bin/"
	@echo "  $(BOLD)make help$(RESET)    — ℹ️  Показать это сообщение"

build:
	@echo "$(BLUE)🛠️  Сборка приложения...$(RESET)"
	@mkdir -p bin
	go build -o $(BIN_PATH) $(CMD_PATH)
	@echo "$(GREEN)✔️  Бинарник собран: $(BIN_PATH)$(RESET)"

run:
	@echo "$(CYAN)🚀 Запуск приложения...$(RESET)"
	@mkdir -p bin
	go run $(CMD_PATH) || true

clean:
	@echo "$(YELLOW)🧹 Очистка bin/ ...$(RESET)"
	rm -rf bin/*
	@echo "$(GREEN)✔️  bin/ очищен$(RESET)"


test:
	@echo "$(BLUE)🧪 Запуск тестов...$(RESET)"
	go test -v ./...

lint:
	@echo "$(YELLOW)🔍 Запуск линтера...$(RESET)"
	golangci-lint run ./...
