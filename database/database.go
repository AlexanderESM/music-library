package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

// db - глобальная переменная для хранения соединения с базой данных.
var (
	db   *gorm.DB
	once sync.Once // Используется для обеспечения одиночного создания подключения
)

// Connect устанавливает соединение с базой данных PostgreSQL.
// Используется шаблон Singleton для предотвращения повторных подключений.
//
// Возвращает *gorm.DB — активное соединение с базой данных.
func Connect() *gorm.DB {
	once.Do(func() {
		// Получаем строку подключения из переменной окружения DATABASE_URL
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("ERROR: DATABASE_URL is not set")
		}

		// Подключаемся к PostgreSQL через GORM
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("ERROR: Failed to connect to the database:", err)
		}

		// Получаем объект sql.DB для настройки соединений
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("ERROR: Failed to get sql.DB from gorm.DB:", err)
		}

		// Настройка пула подключений
		sqlDB.SetMaxOpenConns(10)   // Устанавливаем максимум 10 открытых соединений
		sqlDB.SetMaxIdleConns(5)    // Разрешаем держать до 5 неактивных соединений
		sqlDB.SetConnMaxLifetime(0) // Соединения не будут закрываться автоматически
	})

	return db
}
