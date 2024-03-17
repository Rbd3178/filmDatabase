package models

// FilmRequest
type FilmRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Rating float64 `json:"rating"`
	Actors_IDs  []int  `json:"actors_ids"`
}

// FilmBasic
type FilmBasic struct {
	FilmID int    `json:"film_id" db:"film_id"`
	Title  string `json:"title" db:"title"`
}
