package testdb_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/Rbd3178/filmDatabase/internal/app/store/testdb"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := testdb.New()

	userReq := &models.UserRequest{
		Login:    "JohnDoe",
		Password: "verysecret",
	}
	done, err := s.User().Create(userReq)
	assert.True(t, done)
	assert.NoError(t, err)
}

func TestUserRepository_Find(t *testing.T) {
	s := testdb.New()

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
	s := testdb.New()

	userReq1 := &models.UserRequest{
		Login:    "JohnDoe",
		Password: "verysecret",
	}
	userReq2 := &models.UserRequest{
		Login:    "IvanIvanov",
		Password: "notsosecret",
	}

	s.User().Create(userReq1)
	s.User().Create(userReq2)

	user1, _ := s.User().Find(userReq1.Login)
	user2, _ := s.User().Find(userReq2.Login)
	users, err := s.User().GetAll()
	assert.NoError(t, err)
	assert.Contains(t, users, *user1)
	assert.Contains(t, users, *user2)
}
