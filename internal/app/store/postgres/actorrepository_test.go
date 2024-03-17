package postgres_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestActorRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("actors")

	s := postgres.New(db)

	actorReq := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}

	id, done, err := s.Actor().Create(actorReq)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.NotNil(t, id)
}

func TestActorRepository_Modify(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("actors")

	s := postgres.New(db)

	actorReq := &models.ActorRequest{
		Name:      "Tom Hank",
		Gender:    "male",
		BirthDate: "1959-07-09",
	}

	id, _, _ := s.Actor().Create(actorReq)

	actorReqMod := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "Male",
		BirthDate: "1956-07-09",
	}
	done, err := s.Actor().Modify(id, actorReqMod)
	assert.NoError(t, err)
	assert.True(t, done)

	done, err = s.Actor().Modify(id+10, actorReqMod)
	assert.NoError(t, err)
	assert.False(t, done)
}

func TestActorRepository_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("actors")

	s := postgres.New(db)

	actorReq := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}

	id, _, _ := s.Actor().Create(actorReq)

	done, err := s.Actor().Delete(id)
	assert.NoError(t, err)
	assert.True(t, done)

	done, err = s.Actor().Delete(id)
	assert.NoError(t, err)
	assert.False(t, done)
}

func TestActorRepository_Find(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("films_x_actors, actors, films")

	s := postgres.New(db)

	actorReq := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}

	id, _, _ := s.Actor().Create(actorReq)

	filmReq1 := &models.FilmRequest{
		Title: "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating: 6.8,
		Actors_IDs: []int{id},
	}

	filmReq2 := &models.FilmRequest{
		Title: "Cool title 2",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating: 6.8,
		Actors_IDs: []int{id},
	}

	filmID1, _, _ := s.Film().Create(filmReq1)
	filmID2, _, _ := s.Film().Create(filmReq2)

	actor, done, err := s.Actor().Find(id)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.Contains(t, actor.Films, models.FilmBasic{FilmID: filmID1, Title: filmReq1.Title})
	assert.Contains(t, actor.Films, models.FilmBasic{FilmID: filmID2, Title: filmReq2.Title})
}

func TestActorRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("films_x_actors, actors, films")

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

	filmReq1 := &models.FilmRequest{
		Title: "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating: 6.8,
		Actors_IDs: []int{actorID1},
	}
	filmReq2 := &models.FilmRequest{
		Title: "Cool title 2",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating: 6.8,
		Actors_IDs: []int{actorID1, actorID2},
	}
	filmID1, _, _ := s.Film().Create(filmReq1)
	filmID2, _, _ := s.Film().Create(filmReq2)

	actors, done, err := s.Actor().GetAll()
	assert.NoError(t, err)
	assert.True(t, done)
	for _, actor := range actors {
		if actor.ID == actorID1 {
			assert.Contains(t, actor.Films, models.FilmBasic{FilmID: filmID1, Title: filmReq1.Title})
		}
		assert.Contains(t, actor.Films, models.FilmBasic{FilmID: filmID2, Title: filmReq2.Title})
	}
}