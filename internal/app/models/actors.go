package models

import "time"

// ActorRequest
type ActorRequest struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

// Actor
type Actor struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Gender    string      `json:"gender"`
	BirthDate string      `json:"birth_date"`
	Films     []FilmBasic `json:"films"`
}

// Validate
func (r *ActorRequest) Validate() bool {
	_, err := time.Parse("2006-01-02", r.BirthDate)
	validBirthDate := err == nil || r.BirthDate == ""
	return len(r.Name) >= 1 && len(r.Name) <= 100 && len(r.Gender) <= 20 && validBirthDate
}
