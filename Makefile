# Makefile для Telegram бота "Пидор дня"

.PHONY: build run clean test deps help

# Имя бинарного файла
BINARY_NAME=bot

# Go команды
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Основные цели

help: ## Показать справку
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Установить зависимости
	$(GOMOD) download
	$(GOMOD) tidy

build: deps ## Собрать бота
	$(GOBUILD) -o $(BINARY_NAME) -v .

run: ## Запустить бота
	$(GOCMD) run .

test: ## Запустить тесты
	$(GOTEST) -v ./...

clean: ## Удалить сгенерированные файлы
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# Платформо-специфичная сборка

build-linux: deps ## Собрать для Linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux -v .

build-windows: deps ## Собрать для Windows
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows.exe -v .

build-macos: deps ## Собрать для macOS
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-macos -v .

build-all: build-linux build-windows build-macos ## Собрать для всех платформ

# Разработка

dev: ## Запустить в режиме разработки
	DEBUG=true $(GOCMD) run .

setup: ## Первоначальная настройка проекта
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан файл .env из .env.example"; \
		echo "Не забудьте указать BOT_TOKEN в файле .env"; \
	else \
		echo ".env уже существует"; \
	fi

# Docker (если понадобится в будущем)

docker-build: ## Собрать Docker образ
	docker build -t day-of-the-bot .

docker-run: ## Запустить в Docker
	docker run --env-file .env day-of-the-bot

# Утилиты

fmt: ## Форматировать код
	$(GOCMD) fmt ./...

vet: ## Проверить код
	$(GOCMD) vet ./...

lint: ## Линтер (требует golangci-lint)
	golangci-lint run

check: fmt vet lint test ## Полная проверка кода

install: build ## Установить бота в систему
	sudo mv $(BINARY_NAME) /usr/local/bin/

# По умолчанию
all: clean deps test build