package repository

import (
	"gorm.io/gorm"
	"log"
	"music-library/models"
)

// SongRepository предоставляет методы для работы с песнями в базе данных.
type SongRepository struct {
	DB *gorm.DB // Экземпляр базы данных GORM
}

// SaveSong сохраняет песню в базе данных.
//
// Принимает:
//   - song *models.Song: объект песни для сохранения.
//
// Возвращает:
//   - *models.Song: сохраненный объект песни.
//   - error: ошибка, если сохранение не удалось.
func (repo *SongRepository) SaveSong(song *models.Song) (*models.Song, error) {
	if err := repo.DB.Create(song).Error; err != nil {
		log.Printf("ERROR: Failed to save song. Error: %v\n", err)
		return nil, err
	}
	log.Printf("INFO: Successfully saved song with ID: %d\n", song.ID)
	return song, nil
}

// GetAllSongs получает список всех песен с пагинацией.
//
// Принимает:
//   - page int: номер страницы (начиная с 1).
//   - limit int: количество записей на странице.
//
// Возвращает:
//   - []models.Song: список песен.
//   - error: ошибка, если запрос не удался.
func (repo *SongRepository) GetAllSongs(page int, limit int) ([]models.Song, error) {
	log.Printf("INFO: Retrieving all songs. Page: %d, Limit: %d\n", page, limit)
	var songs []models.Song
	offset := (page - 1) * limit

	if err := repo.DB.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		log.Printf("ERROR: Failed to retrieve songs. Page: %d, Limit: %d, Error: %v\n", page, limit, err)
		return nil, err
	}
	log.Printf("INFO: Successfully retrieved %d songs.\n", len(songs))
	return songs, nil
}

// GetSongByID получает песню по ее уникальному идентификатору.
//
// Принимает:
//   - id uint: ID песни.
//
// Возвращает:
//   - *models.Song: найденная песня.
//   - error: ошибка, если песня не найдена.
func (repo *SongRepository) GetSongByID(id uint) (*models.Song, error) {
	log.Printf("INFO: Retrieving song with ID: %d\n", id)
	var song models.Song
	if err := repo.DB.First(&song, id).Error; err != nil {
		log.Printf("ERROR: Failed to retrieve song with ID: %d, Error: %v\n", id, err)
		return nil, err
	}
	log.Printf("INFO: Successfully retrieved song with ID: %d\n", id)
	return &song, nil
}

// UpdateSong обновляет существующую песню в базе данных.
//
// Принимает:
//   - song *models.Song: обновленный объект песни.
//
// Возвращает:
//   - *models.Song: обновленный объект песни.
//   - error: ошибка, если обновление не удалось.
func (repo *SongRepository) UpdateSong(song *models.Song) (*models.Song, error) {
	log.Printf("INFO: Updating song with ID: %d\n", song.ID)
	if err := repo.DB.Save(song).Error; err != nil {
		log.Printf("ERROR: Failed to update song with ID: %d, Error: %v\n", song.ID, err)
		return nil, err
	}
	log.Printf("INFO: Successfully updated song with ID: %d\n", song.ID)
	return song, nil
}

// DeleteSong удаляет песню по ее уникальному идентификатору.
//
// Принимает:
//   - id uint: ID песни.
//
// Возвращает:
//   - error: ошибка, если удаление не удалось.
func (repo *SongRepository) DeleteSong(id uint) error {
	log.Printf("INFO: Deleting song with ID: %d\n", id)
	if err := repo.DB.Delete(&models.Song{}, id).Error; err != nil {
		log.Printf("ERROR: Failed to delete song with ID: %d, Error: %v\n", id, err)
		return err
	}
	log.Printf("INFO: Successfully deleted song with ID: %d\n", id)
	return nil
}
