package testdb

import (
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
)

// FilmRepository
type FilmRepository struct {
	store  *Store
	films  map[int]*models.Film
	actors map[int]*models.ActorBasic
}

// Create
func (r *FilmRepository) Create(f *models.FilmRequest) (int, error) {
	id := len(r.films) + 1
	film := models.Film{
		Title:       f.Title,
		Description: f.Description,
		ReleaseDate: f.ReleaseDate,
		Rating:      f.Rating,
	}
	for _, actorID := range f.ActorsIDs {
		actor, ok := r.actors[actorID]
		if !ok {
			continue
		}
		film.Actors = append(film.Actors, *actor)
	}
	r.films[id] = &film
	return id, nil
}

// GetAll
func (r *FilmRepository) GetAll(orderBy string, order string, searchTitle string, searchActor string) ([]models.Film, error) {
	films := make([]models.Film, 0)
	for _, film := range r.films {
		films = append(films, *film)
	}
	return films, nil
}

// Delete
func (r *FilmRepository) Delete(id int) (bool, error) {
	_, ok := r.films[id]
	if !ok {
		return false, nil
	}

	delete(r.films, id)

	return true, nil
}

// Find
func (r *FilmRepository) Find(id int) (*models.Film, error) {
	film, ok := r.films[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return film, nil
}

// Modify
func (r *FilmRepository) Modify(id int, f *models.FilmRequest) (bool, error) {
	_, ok := r.films[id]
	if !ok {
		return false, nil
	}

	r.films[id].Title = f.Title
	r.films[id].Description = f.Description
	r.films[id].ReleaseDate = f.ReleaseDate
	r.films[id].Rating = f.Rating

	r.films[id].Actors = make([]models.ActorBasic, 0)
	for _, actorID := range f.ActorsIDs {
		actor, ok := r.actors[actorID]
		if !ok {
			continue
		}
		r.films[id].Actors = append(r.films[id].Actors, *actor)
	}

	return true, nil
}
