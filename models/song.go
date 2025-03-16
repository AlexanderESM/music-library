package models

import "time"

// Song представляет собой модель песни.
type Song struct {
	ID          int       `json:"id"`
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	Text        string    `json:"text"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

// SongDetail представляет более подробную информацию о песне.
type SongDetail struct {
	Link        string `json:"link"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
}
