package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"music-library/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// SongEnrichment структура для данных, обогащающих информацию о песне
type SongEnrichment struct {
	Group       string `json:"group"`        // Группа исполнителей
	Song        string `json:"song"`         // Название песни
	ReleaseDate string `json:"release_date"` // Дата релиза песни
	Text        string `json:"text"`         // Текст песни
	Link        string `json:"link"`         // Ссылка на внешний источник
}

// GetSongInfo обрабатывает запросы для получения информации о песне и добавляет её в базу данных
func GetSongInfo(c *gin.Context) {
	group := c.Query("group") // Получаем название группы из параметров запроса
	song := c.Query("song")   // Получаем название песни из параметров запроса

	// Проверка на наличие обязательных параметров
	if group == "" || song == "" {
		log.Println("ERROR: Missing 'group' or 'song' query parameters")
		c.String(http.StatusBadRequest, "bad request: missing required parameters")
		return
	}

	var songRecord models.Song
	// Используем db из контекста (db передается в main)
	db := c.MustGet("db").(*gorm.DB)

	// Поиск песни в базе данных
	err := db.Where("\"group\" = ? AND song = ?", group, song).First(&songRecord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("INFO: Song '%s' by '%s' not found in database.", song, group)

		// Запрос к внешнему API для получения информации о песне
		songDetail, shouldReturn := GetSongDetailFromAPI(group, song, c)
		if shouldReturn {
			return
		}

		// Конвертируем ReleaseDate из строки в time.Time
		releaseDate, err := time.Parse("2006-01-02", songDetail.ReleaseDate) // Подгоните формат под ваш случай
		if err != nil {
			log.Printf("ERROR: Failed to parse release date: %v", err)
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		// Создание новой записи о песне
		newSong := models.Song{
			Group:       group,
			Song:        song,
			ReleaseDate: releaseDate, // Присваиваем time.Time
			Text:        songDetail.Text,
			Link:        songDetail.Link,
		}

		// Сохраняем новую песню в базу данных
		if err := db.Create(&newSong).Error; err != nil {
			log.Printf("ERROR: Failed to add new song to the database: %v", err)
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		log.Printf("INFO: Added new song to the database: %v", newSong)
		songRecord = newSong
	} else if err != nil {
		log.Printf("ERROR: Database error: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	// Формируем ответ с деталями песни
	songDetail := models.SongDetail{
		ReleaseDate: songRecord.ReleaseDate.Format("2006-01-02"), // Форматируем в строку для ответа
		Text:        songRecord.Text,
		Link:        songRecord.Link,
	}

	// Дополнительное обогащение данных из файла
	enrichSongFromJSON(&songDetail, group, song)
	c.JSON(http.StatusOK, songDetail)
}

// GetSongs возвращает список всех песен из базы данных
func GetSongs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var songs []models.Song

	// Получаем все песни из базы данных
	if err := db.Find(&songs).Error; err != nil {
		log.Printf("ERROR: Failed to fetch songs: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	// Отправляем список песен в ответе
	c.JSON(http.StatusOK, songs)
}

// GetSongTextWithPagination возвращает текст песни с пагинацией
func GetSongTextWithPagination(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	songID := c.Param("id") // Получаем ID песни из параметров маршрута

	// Получаем песню по ID
	var song models.Song
	if err := db.First(&song, songID).Error; err != nil {
		log.Printf("ERROR: Song not found with ID: %v", songID)
		c.String(http.StatusNotFound, "Song not found")
		return
	}

	// Разбиваем текст песни на строки или части (например, по строкам)
	verses := splitTextIntoParts(song.Text)

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Получаем подмножество текста с учетом пагинации
	start := (page - 1) * limit
	end := start + limit
	if end > len(verses) {
		end = len(verses)
	}

	// Возвращаем выбранные части текста
	c.JSON(http.StatusOK, verses[start:end])
}

// splitTextIntoParts разделяет текст на части (например, по строкам)
func splitTextIntoParts(text string) []string {
	// Разделяем текст по строкам
	return strings.Split(text, "\n")
}

// GetSongDetailFromJSON получает информацию о песне из JSON файла
func GetSongDetailFromJSON(group, song string) (models.SongDetail, error) {
	// Открываем JSON файл
	jsonFile, err := os.Open("song_enrichment.json")
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("не удалось открыть JSON файл: %v", err)
	}
	defer jsonFile.Close()

	// Читаем данные из файла
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("не удалось прочитать JSON файл: %v", err)
	}

	var enrichmentData SongEnrichment
	if err := json.Unmarshal(byteValue, &enrichmentData); err != nil {
		return models.SongDetail{}, fmt.Errorf("не удалось разобрать JSON файл: %v", err)
	}

	// Если группа и песня совпадают, возвращаем данные
	if enrichmentData.Group == group && enrichmentData.Song == song {
		return models.SongDetail{
			ReleaseDate: enrichmentData.ReleaseDate,
			Text:        enrichmentData.Text,
			Link:        enrichmentData.Link,
		}, nil
	}

	return models.SongDetail{}, fmt.Errorf("песня не найдена в JSON файле")
}

// GetSongDetailFromAPI выполняет запрос к внешнему API
func GetSongDetailFromAPI(group, song string, c *gin.Context) (models.SongDetail, bool) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	apiURL := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", encodedGroup, encodedSong)

	// Запрос к внешнему API
	response, err := http.Get(apiURL)
	if err != nil {
		log.Printf("ERROR: Failed to request external API: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return models.SongDetail{}, true
	}
	defer response.Body.Close()

	// Проверка статуса ответа
	if response.StatusCode != http.StatusOK {
		log.Printf("WARNING: External API returned status code %d", response.StatusCode)
		c.String(http.StatusInternalServerError, "failed to retrieve song details from external API")
		return models.SongDetail{}, true
	}

	var apiData models.SongDetail
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read API response: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return models.SongDetail{}, true
	}

	if err := json.Unmarshal(body, &apiData); err != nil {
		log.Printf("ERROR: Failed to parse API response: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return models.SongDetail{}, true
	}

	return apiData, false
}

// enrichSongFromJSON читает данные из JSON-файла и обогащает информацию о песне
func enrichSongFromJSON(songDetail *models.SongDetail, group, song string) {
	jsonFile, err := os.Open("song_enrichment.json")
	if err != nil {
		log.Printf("ERROR: Could not open JSON file: %v", err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Printf("ERROR: Failed to read JSON file: %v", err)
		return
	}

	var enrichmentData SongEnrichment
	if err := json.Unmarshal(byteValue, &enrichmentData); err != nil {
		log.Printf("ERROR: Could not parse JSON file: %v", err)
		return
	}

	// Обогащение информации о песне
	if enrichmentData.Group == group && enrichmentData.Song == song {
		// Преобразуем release date в time.Time
		releaseDate, err := time.Parse("2006-01-02", enrichmentData.ReleaseDate)
		if err != nil {
			log.Printf("ERROR: Failed to parse release date from enrichment data: %v", err)
			return
		}
		// Присваиваем ReleaseDate как time.Time в songDetail
		songDetail.ReleaseDate = releaseDate.Format("2006-01-02") // Преобразуем в строку для SongDetail
		songDetail.Text = enrichmentData.Text
		songDetail.Link = enrichmentData.Link
	}
}

// UpdateSong обновляет песню
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var song models.Song
	if err := db.First(&song, id).Error; err != nil {
		log.Printf("ERROR: Song with ID %s not found", id)
		c.String(http.StatusNotFound, "not found")
		return
	}

	// Обработка JSON данных из запроса
	if err := c.ShouldBindJSON(&song); err != nil {
		log.Printf("ERROR: Invalid song data: %v", err)
		c.String(http.StatusBadRequest, "invalid input")
		return
	}

	// Обновление записи в базе данных
	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&song).Error
	}); err != nil {
		log.Printf("ERROR: Failed to update song with ID %s: %v", id, err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	log.Printf("INFO: Updated song with ID %s", id)
	c.JSON(http.StatusOK, song)
}

// DeleteSong удаляет песню
func DeleteSong(c *gin.Context) {
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Delete(&models.Song{}, id).Error; err != nil {
		log.Printf("ERROR: Failed to delete song with ID %s: %v", id, err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	log.Printf("INFO: Deleted song with ID %s", id)
	c.JSON(http.StatusOK, map[string]interface{}{"id #" + id: "deleted"})
}
