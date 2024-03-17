package postgres_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestFilmRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("films_x_actors, films, actors")

	s := postgres.New(db)

	actorReq1 := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}

	actorReq2 := &models.ActorRequest{
		Name:      "Sophie Patel",
		Gender:    "female",
		BirthDate: "1997-11-03",
	}

	actorID1, _, _ := s.Actor().Create(actorReq1)
	actorID2, _, _ := s.Actor().Create(actorReq2)

	filmReq := &models.FilmRequest{
		Title: "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating: 6.8,
		Actors_IDs: []int{actorID1, actorID2, actorID2 + 10},
	}

	id, done, err := s.Film().Create(filmReq)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.NotNil(t, id)
}
