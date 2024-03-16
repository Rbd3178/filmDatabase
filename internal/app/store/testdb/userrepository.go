package testdb

import (
	"errors"

	"github.com/Rbd3178/filmDatabase/internal/app/hasher"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
)

// UserRepository
type UserRepository struct {
	store *Store
	users map[string]*models.User
}

// Create
func (r *UserRepository) Create(u *models.UserRequest) error {
	hashedPassword, err := hasher.HashPassword(u.Password)
	if err != nil {
		return err
	}
	if _, ok := r.users[u.Login]; ok {
		return errors.New("login is taken")
	}
	user := &models.User{
		Login: u.Login,
		HashedPassword: hashedPassword,
		IsAdmin: false,
	}
	r.users[u.Login] = user

	return nil
}

// Find ...
func (r *UserRepository) Find(login string) (*models.User, error) {
	u, ok := r.users[login]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

// GetAll
func (r *UserRepository) GetAll() ([]models.User, error) {
	users := make([]models.User, 0) 
	for _, u := range r.users {
		users = append(users, *u)
	}
	return users, nil
}
