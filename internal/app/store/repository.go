package store

import (
	//"github.com/Rbd3178/filmDatabase/internal/app/entities"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
)

// UserRepository
type UserRepository interface {
	Create(*models.UserRequest) error
	Find(string) (*models.User, error)
	GetAll() ([]models.User, error)
}

// FilmRepository
type FilmRepository interface {
	Create(*models.FilmRequest) (int, bool, error)
}

// ActorRepository
type ActorRepository interface {
	Create(*models.ActorRequest) (int, bool, error)
	Modify(int, *models.ActorRequest) (bool, error)
	Delete(int) (bool, error)
	Find(int) (*models.Actor, bool, error)
	GetAll() ([]models.Actor, bool, error)
}