package models

import "time"

// FilmRequest
type FilmRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"release_date"`
	Rating      float64 `json:"rating"`
	Actors_IDs  []int   `json:"actors_ids"`
}

// FilmBasic
type FilmBasic struct {
	FilmID int    `json:"film_id" db:"film_id"`
	Title  string `json:"title" db:"title"`
}

// VerifyForInsert
func (r *FilmRequest) ValidateForInsert() bool {
	_, err := time.Parse("2006-01-02", r.ReleaseDate)
	validReleaseDate := err == nil
	validTitle := len(r.Title) >= 1 && len(r.Title) <= 150
	validDescription := len(r.Description) <= 1000
	validRating := r.Rating >= 0 && r.Rating <= 10
	return validTitle && validDescription && validReleaseDate && validRating
}
