package testdb

import (
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
)

// Store
type Store struct {
	userRepository *UserRepository
	actorRepository *ActorRepository
	filmRepository *FilmRepository
}

// New
func New() *Store {
	return &Store{}
}

// User
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]*models.User),
	}

	return s.userRepository
}

// Actor
func (s *Store) Actor() store.ActorRepository {
	if s.actorRepository != nil {
		return s.actorRepository
	}

	s.actorRepository = &ActorRepository{
		store: s,
		actors: make(map[int]*models.Actor),
	}

	return s.actorRepository
}

// Film
func (s *Store) Film() store.FilmRepository {
	if s.filmRepository != nil {
		return s.filmRepository
	}

	s.filmRepository = &FilmRepository{
		store: s,
		films: make(map[int]*models.Film),
		actors: make(map[int]*models.ActorBasic),
	}
	s.filmRepository.actors[1] = &models.ActorBasic{
		ActorID: 1,
		Name: "First Actor",
	}
	s.filmRepository.actors[2] = &models.ActorBasic{
		ActorID: 2,
		Name: "Second Actor",
	}
	s.filmRepository.actors[3] = &models.ActorBasic{
		ActorID: 3,
		Name: "Third Actor",
	}
	return s.filmRepository
}