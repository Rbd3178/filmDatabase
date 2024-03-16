package postgres

import (
	"github.com/Rbd3178/filmDatabase/internal/app/entities"
	"github.com/Rbd3178/filmDatabase/internal/app/hasher"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
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
	err = r.store.db.QueryRow(
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

// FindByLogin
func (r *UserRepository) FindByLogin(login string) (*entities.User, error) {
	u := &entities.User{}
	err := r.store.db.Select(u, "SELECT login, hashed_password, is_admin FROM users WHERE login = $1", &login)
	if err != nil {
		return nil, err
	}
	return u, nil
}
