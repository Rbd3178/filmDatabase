package store

import (
	"github.com/Rbd3178/filmDatabase/internal/app/hasher"
	"github.com/Rbd3178/filmDatabase/internal/app/model"
)

// UserRepository
type UserRepository struct {
	store *Store
}

// Create
func (r *UserRepository) Create(u *model.UserRequest) (*model.User, error) {
	hashedPassword, err := hasher.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	err = r.store.db.QueryRow(
		"INSERT INTO users (login, hashed_password, is_admin) VALUES ($1, $2, $3)",
		u.Login,
		hashedPassword,
		false,
	).Scan()
	if err != nil {
		return nil, err
	}

	return &model.User{Login: u.Login, HashedPassword: hashedPassword, Role: "user"}, nil
}

// FindByLogin
func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	return nil, nil
}
