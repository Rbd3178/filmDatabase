package models

// ActorRequest
type ActorRequest struct {
	Name string `json:"name"`
	Gender string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

// Actor
type Actor struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Gender string `json:"gender"`
	BirthDate string `json:"birth_date"`
	Films []FilmBasic `json:"films"`
}