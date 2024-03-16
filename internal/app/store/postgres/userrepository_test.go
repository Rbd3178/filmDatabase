package postgres_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/Rbd3178/filmDatabase/internal/app/store/postgres"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("users")

	s := postgres.New(db)
	userReq := &models.UserRequest{
		Login:    "JohnDoe",
		Password: "verysecret",
	}
	assert.NoError(t, s.User().Create(userReq))
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("users")

	s := postgres.New(db)
	userReq := &models.UserRequest{
		Login:    "JohnDoe",
		Password: "verysecret",
	}
	_, err := s.User().Find(userReq.Login)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(userReq)
	u, err := s.User().Find(userReq.Login)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("users")

	s := postgres.New(db)
	userReq := &models.UserRequest{
		Login:    "JohnDoe",
		Password: "verysecret",
	}
	_, err := s.User().Find(userReq.Login)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(userReq)
	u, err := s.User().Find(userReq.Login)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
