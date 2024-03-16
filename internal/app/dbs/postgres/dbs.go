package postgres

import (
	"github.com/Rbd3178/filmDatabase/internal/app/dbs"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// Store
type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

// New
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// User
func (s *Store) User() dbs.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
