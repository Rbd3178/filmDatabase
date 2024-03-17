package postgres

import (
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// Store
type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
	filmRepository *FilmRepository
	actorRepository *ActorRepository
}

// New
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// User
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

// Film
func (s *Store) Film() store.FilmRepository {
	if s.filmRepository != nil {
		return s.filmRepository
	}

	s.filmRepository = &FilmRepository{
		store: s,
	}

	return s.filmRepository
}

// Actor
func (s *Store) Actor() store.ActorRepository {
	if s.actorRepository != nil {
		return s.actorRepository
	}

	s.actorRepository = &ActorRepository{
		store: s,
	}

	return s.actorRepository
}