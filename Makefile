.PHONY: build run dev clean test deps

# Сборка приложения
build:
	go build -o bin/telegram-bot cmd/main.go

# Запуск в продакшене
run: build
	./bin/telegram-bot

# Запуск в режиме разработки
dev:
	go run cmd/main.go

# Очистка артефактов сборки
clean:
	rm -rf bin/

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Тестирование
test:
	go test ./...

# Проверка форматирования кода
fmt:
	go fmt ./...

# Линтер
lint:
	golangci-lint run

# Создание директорий для веб-файлов
setup-web:
	mkdir -p web/static

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build      - Сборка приложения"
	@echo "  run        - Запуск в продакшене"
	@echo "  dev        - Запуск в режиме разработки"
	@echo "  clean      - Очистка артефактов сборки"
	@echo "  deps       - Установка зависимостей"
	@echo "  test       - Запуск тестов"
	@echo "  fmt        - Форматирование кода"
	@echo "  lint       - Проверка кода линтером"
	@echo "  setup-web  - Создание директорий для веб-файлов"
	@echo "  help       - Показать эту справку" 