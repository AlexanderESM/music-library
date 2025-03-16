package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"music-library/models"
)

// Migrate выполняет автоматическую миграцию моделей в базе данных.
// Использует GORM AutoMigrate для создания или обновления таблиц.
//
// Аргументы:
//   - db: *gorm.DB — соединение с базой данных.
//
// Возвращает:
//   - error: Ошибка миграции (если произошла), иначе nil.
func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("ERROR: database connection is nil")
	}

	// Определяем список моделей для миграции
	modelsToMigrate := []interface{}{
		&models.Song{}, // Таблица для хранения информации о песнях
		// Добавьте другие модели сюда при необходимости
	}

	// Выполняем миграцию
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		log.Printf("ERROR: Migration failed: %v", err)
		return err
	}

	log.Println("INFO: Migration completed successfully")
	return nil
}
