package config

import (
	"github.com/joho/godotenv"
	"log"
)

// LoadEnv загружает переменные окружения из .env файла в текущий процесс.
// Используется для конфигурирования приложения без необходимости жестко прописывать параметры в коде.
//
// Если файл .env отсутствует или не может быть загружен, программа завершает работу с фатальной ошибкой.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: Failed to load .env file. Ensure the file exists and is readable.")
	}

	log.Println("INFO: .env file successfully loaded.")
}
