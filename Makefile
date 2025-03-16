# Указывает, что цели не являются файлами и всегда должны выполняться заново
.PHONY: run swag-generate all

# Основная цель: выполняет генерацию Swagger-документации и запускает приложение
all: swag-generate run

# Запуск приложения
# Эта команда выполняет `go run cmd/main.go`, который запускает основное приложение.
run:
	go run cmd/main.go

# Генерация Swagger-документации с помощью swaggo/swag
# - `cd cmd` — переходит в каталог `cmd`
# - `swag init` — создает Swagger-документацию
# - `-g ../cmd/main.go` — указывает главный файл приложения
# - `-d ../config,../models,../controllers,../database,../repository` — определяет каталоги, используемые для генерации документации
# - `-o ../docs` — задает папку для сохранения документации
swag-generate:
	cd cmd && swag init -g ../cmd/main.go -d ../config,../models,../controllers,../database,../repository -o ../docs
