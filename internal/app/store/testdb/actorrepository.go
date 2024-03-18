package testdb

import (
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
)

// ActorRepository
type ActorRepository struct {
	store  *Store
	actors map[int]*models.Actor
}

var (
	sampleFilm1 = models.FilmBasic{
		FilmID: 1,
		Title: "Title One",
	}

	sampleFilm2 = models.FilmBasic{
		FilmID: 2,
		Title: "Title Two",
	}
)

// Create
func (r *ActorRepository) Create(a *models.ActorRequest) (int, error) {
	id := len(r.actors) + 1
	actor := models.Actor{
		Name: a.Name,
		Gender: a.Gender,
		BirthDate: a.BirthDate,
		Films: []models.FilmBasic{sampleFilm1, sampleFilm2},
	}
	r.actors[id] = &actor
	return id, nil
}

// Modify
func (r *ActorRepository) Modify(id int, a *models.ActorRequest) (bool, error) {
	_, ok := r.actors[id]
	if !ok {
		return false, nil
	}

	r.actors[id].Name = a.Name
	r.actors[id].Gender = a.Gender
	r.actors[id].BirthDate = a.BirthDate

	return true, nil
}

// Delete
func (r *ActorRepository) Delete(id int) (bool, error) {
	_, ok := r.actors[id]
	if !ok {
		return false, nil
	}

	delete(r.actors, id)

	return true, nil
}

// Find
func (r *ActorRepository) Find(id int) (*models.Actor, error) {
	actor, ok := r.actors[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return actor, nil
}

// GetAll
func (r *ActorRepository) GetAll() ([]models.Actor, error) {
	actors := make([]models.Actor, 0, len(r.actors))
	for _, actor := range r.actors {
		actors = append(actors, *actor)
	}

	return actors, nil
}
