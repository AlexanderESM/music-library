package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"music-library/config"
	"music-library/controllers"
	"music-library/database"
	_ "music-library/docs"
	"net/http"
)

// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен.
// @host localhost:8080
// @BasePath /

func main() {
	// 1. Загрузка переменных окружения из .env файла
	config.LoadEnv()
	log.Println("INFO: Environment variables loaded.")

	// 2. Подключение к базе данных
	db := database.Connect()
	log.Println("INFO: Database connection established.")

	// 3. Выполнение миграций базы данных (создание таблиц, если их нет)
	database.Migrate(db)
	log.Println("INFO: Database migrations completed.")

	// 4. Инициализация HTTP-сервера с помощью Gin
	router := gin.Default()

	// 5. Определение маршрутов основного API
	router.GET("/info", controllers.GetSongInfo)                           // Получение информации о песне
	router.GET("/songs", controllers.GetSongs)                             // Получение списка всех песен
	router.GET("/songs/:id/verses", controllers.GetSongTextWithPagination) // Получение текста песни с пагинацией
	router.PUT("/songs/:id", controllers.UpdateSong)                       // Обновление информации о песне по ID
	router.DELETE("/songs/:id", controllers.DeleteSong)                    // Удаление песни по ID

	// 6. Swagger-документация доступна по адресу http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("INFO: Swagger documentation is available at http://localhost:8080/swagger/index.html")

	// 7. Запуск второго сервера на порту 8081 для эмуляции внешнего API
	go startMockServer()

	// 8. Запуск основного сервера на порту 8080
	log.Println("INFO: Starting the main server on port 8080...")
	log.Fatal(router.Run(":8080")) // Запуск основного HTTP-сервера и логирование фатальных ошибок
}

// startMockServer запускает тестовый сервер на порту 8081 для эмуляции внешнего API
func startMockServer() {
	testRouter := gin.Default()

	// 1. Определяем маршрут для получения информации о песне из внешнего JSON
	testRouter.GET("/info", func(c *gin.Context) {
		group := c.Query("group") // Получаем параметр "group" из URL
		song := c.Query("song")   // Получаем параметр "song" из URL

		// 2. Проверка параметров запроса
		if group == "" || song == "" {
			log.Println("DEBUG: Missing request parameters: group or song.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
			return
		}

		// 3. Получение информации о песне из JSON
		songDetail, err := controllers.GetSongDetailFromJSON(group, song)
		if err != nil {
			log.Printf("DEBUG: Error fetching song details: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
			return
		}

		// 4. Успешный ответ с информацией о песне
		log.Printf("INFO: Request to /info succeeded for group: %s, song: %s\n", group, song)
		c.JSON(http.StatusOK, songDetail)
	})

	// 5. Запуск тестового сервера на порту 8081
	if err := testRouter.Run(":8081"); err != nil {
		log.Fatalf("ERROR: Failed to start the test server: %v", err)
	}
}
