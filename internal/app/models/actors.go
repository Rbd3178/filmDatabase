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

// ValidateForInsert
func (r *ActorRequest) ValidateForInsert() bool {
	_, err := time.Parse("2006-01-02", r.BirthDate)
	validBirthDate := err == nil
	validName := len(r.Name) >= 1 && len(r.Name) <= 100
	validGender := len(r.Gender) <= 20
	return validName && validGender && validBirthDate
}

// ValidateForUpdate
func (r *ActorRequest) ValidateForUpdate() bool {
	_, err := time.Parse("2006-01-02", r.BirthDate)
	validBirthDate := err == nil || r.BirthDate == ""
	validName := len(r.Name) <= 100
	validGender := len(r.Gender) <= 20
	return validName && validGender && validBirthDate
}
