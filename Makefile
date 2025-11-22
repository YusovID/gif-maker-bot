# ====================================================================================
# CONFIGURATION
# ====================================================================================

APP_NAME := gif-maker-bot
COMPOSE_FILE := compose.yml

-include .env

export

# ====================================================================================
# COMMANDS & TOOLS
# ====================================================================================

RM := rm -f
SLEEP := sleep

COMPOSE := docker compose -f $(COMPOSE_FILE)

GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint

# ====================================================================================
# SETUP
# ====================================================================================

.DEFAULT_GOAL := help

.PHONY: all help build up start stop restart down nuke logs ps clean fmt lint test test-cover tools

# ====================================================================================
# GENERAL COMMANDS
# ====================================================================================

all: fmt lint test ## Запустить форматирование, линтер и тесты

help: ## Показать этот список команд и их описания
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "; printf "  \033[36m%-20s\033[0m %s\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?## / { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST) | sort

# ====================================================================================
# DOCKER COMPOSE MANAGEMENT
# ====================================================================================

build: ## Собрать или пересобрать образы сервисов
	@echo "Building service images..."
	@$(COMPOSE) build

up: down build ## Собрать образы и запустить сервисы в фоне
	@echo "Starting services..."
	@$(COMPOSE) up -d

start: ## Запустить ранее остановленные контейнеры
	@echo "Starting existing containers..."
	@$(COMPOSE) start

stop: ## Остановить запущенные сервисы
	@echo "Stopping services..."
	@$(COMPOSE) stop

restart: stop start ## Перезапустить сервисы

down: ## Остановить и удалить контейнеры/сети (тома сохраняются)
	@echo "Tearing down services..."
	@$(COMPOSE) down --remove-orphans

nuke: ## ВНИМАНИЕ: Полностью удалить всё (контейнеры, сети, ТОМА)
	@echo "Nuking the entire environment (containers, networks, VOLUMES)..."
	@$(COMPOSE) down -v --remove-orphans

logs: ## Показать логи всех сервисов в реальном времени
	@$(COMPOSE) logs -f

ps: ## Показать статус запущенных контейнеров
	@$(COMPOSE) ps

# ====================================================================================
# GO BUILD & TEST
# ====================================================================================

fmt: ## Отформатировать весь Go код
	@echo "Formatting Go files..."
	@gofmt -w .

lint: tools ## Запустить линтер для проверки качества кода
	@echo "Running linter..."
	@$(GOLANGCI_LINT) run ./...

test: ## Запустить unit-тесты (без интеграционных)
	@echo "Running fast tests..."
	@go test -v -race -short ./...

test-cover: ## Запустить ВСЕ тесты и сгенерировать HTML-отчет о покрытии
	@echo "Running all tests with coverage..."
	@go test -race -short -coverprofile=unit.cover ./...
	@go test -race -tags=integration -coverprofile=integration.cover ./...

	@echo "Merging coverage profiles..."
	@echo "mode: set" > coverage.out
	@cat unit.cover integration.cover | grep -v "^mode:" >> coverage.out

	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html

	@echo "Cleaning up intermediate files..."
	@$(RM) unit.cover integration.cover coverage.out

	@echo "Coverage report successfully generated: open coverage.html"

clean: ## Очистить все артефакты сборки и тестирования
	@echo "Cleaning up build and test artifacts..."
	@$(RM) coverage.html coverage.out unit.cover integration.cover
	@$(RM) *.test *.exe

tools: ## Установить/обновить зависимости для утилит
	@echo "Syncing tools dependencies..."
	@go mod -C tools tidy