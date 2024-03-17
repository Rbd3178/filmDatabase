package store

import (
	"github.com/Rbd3178/filmDatabase/internal/app/models"
)

// UserRepository
type UserRepository interface {
	Create(*models.UserRequest) (bool, error)
	Find(string) (*models.User, error)
	GetAll() ([]models.User, error)
}

// FilmRepository
type FilmRepository interface {
	Create(*models.FilmRequest) (int, bool, error)
	GetAll(string, string, string) ([]models.Film, error)
}

// ActorRepository
type ActorRepository interface {
	Create(*models.ActorRequest) (int, error)
	Modify(int, *models.ActorRequest) (bool, error)
	Delete(int) (bool, error)
	Find(int) (*models.Actor, error)
	GetAll() ([]models.Actor, error)
}