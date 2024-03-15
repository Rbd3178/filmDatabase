package store

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// Store
type Store struct {
	config *Config
	db     *sqlx.DB
}

// New
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open
func (s *Store) Open() error {
	db, err := sqlx.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

// Close
func (s *Store) Close() {
	s.db.Close()
}
