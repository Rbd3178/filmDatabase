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
