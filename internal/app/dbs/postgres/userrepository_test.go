package postgres_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/dbs/postgres"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("users")

	s := postgres.New(db)
	u := &models.UserRequest{
		Login: "JohnDoe",
		Password: "verysecret",
	}
	err := s.User().Create(u)
	if err != nil {
		t.Fatal(err)
	}
}
