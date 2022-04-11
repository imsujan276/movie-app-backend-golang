package models

import (
	"database/sql"
	"time"
)

// Models is the wrapper for Database
type Models struct {
	DB DBModel
}

// NewModels returns models with DB pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// Movie is the type for movie
type Movie struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Year        int            `json:"year"`
	ReleaseDate time.Time      `json:"release_date"`
	RunTime     int            `json:"run_time"`
	Rating      int            `json:"rating"`
	MPAARating  string         `json:"mpaa_rating"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	MovieGenre  map[int]string `json:"genres"`
	Poster      string         `json:"poster"`
}

// Genre is the type for Genre
type Genre struct {
	ID        int       `json:"id"`
	GenreName string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// MovieGenre is the type for MovieGenre
type MovieGenre struct {
	ID        int       `json:"-"`
	GenreID   int       `json:"genre_id"`
	MovieID   int       `json:"-"`
	Genre     Genre     `json:"genre"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// User is the type for User
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
