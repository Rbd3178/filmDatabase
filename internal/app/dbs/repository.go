package dbs

import (
	"github.com/Rbd3178/filmDatabase/internal/app/entities"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
)

// UserRepository
type UserRepository interface {
	Create(*models.UserRequest) error
	FindByLogin(string) (*entities.User, error)
}
