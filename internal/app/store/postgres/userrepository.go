package postgres

import (
	"database/sql"

	"github.com/Rbd3178/filmDatabase/internal/app/hasher"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	//"github.com/jmoiron/sqlx"
)

// UserRepository
type UserRepository struct {
	store *Store
}

// Create
func (r *UserRepository) Create(u *models.UserRequest) error {
	hashedPassword, err := hasher.HashPassword(u.Password)
	if err != nil {
		return err
	}
	err = r.store.db.QueryRowx(
		"INSERT INTO users (login, hashed_password, is_admin) VALUES ($1, $2, $3) RETURNING login",
		u.Login,
		hashedPassword,
		false,
	).Scan(&u.Login)
	if err != nil {
		return err
	}

	return nil
}

// Find
func (r *UserRepository) Find(login string) (*models.User, error) {
	u := &models.User{}
	err := r.store.db.Get(
		u,
		"SELECT * FROM users WHERE login = $1",
		&login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// GetAll
func (r *UserRepository) GetAll() ([]models.User, error) {
	users := make([]models.User, 0)
	err := r.store.db.Select(
		&users,
		"SELECT * FROM users",
	)
	if err != nil {
		return nil, err
	}
	return users, nil
}
