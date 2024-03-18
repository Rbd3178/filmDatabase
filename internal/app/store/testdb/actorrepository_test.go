package testdb_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/Rbd3178/filmDatabase/internal/app/store/testdb"
	"github.com/stretchr/testify/assert"
)

var (
	sampleFilm1 = models.FilmBasic{
		FilmID: 1,
		Title: "Title One",
	}

	sampleFilm2 = models.FilmBasic{
		FilmID: 2,
		Title: "Title Two",
	}
)

func TestActorRepository_Create(t *testing.T) {
	s := testdb.New()

	actorReq := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}

	id, err := s.Actor().Create(actorReq)
	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestActorRepository_Modify(t *testing.T) {
	s := testdb.New()

	actorReq := &models.ActorRequest{
		Name:      "Tom Hank",
		Gender:    "male",
		BirthDate: "1959-07-09",
	}

	id, _ := s.Actor().Create(actorReq)

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
	s := testdb.New()

	actorReq := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}
	id, _ := s.Actor().Create(actorReq)

	filmReq := &models.FilmRequest{
		Title:       "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      6.8,
		ActorsIDs:   []int{id},
	}
	s.Film().Create(filmReq)

	done, err := s.Actor().Delete(id)
	assert.NoError(t, err)
	assert.True(t, done)

	done, err = s.Actor().Delete(id)
	assert.NoError(t, err)
	assert.False(t, done)
}

func TestActorRepository_Find(t *testing.T) {
	s := testdb.New()

	actorReq := &models.ActorRequest{
		Name:      "Tom Hanks",
		Gender:    "male",
		BirthDate: "1956-07-09",
	}

	id, _ := s.Actor().Create(actorReq)

	actor, err := s.Actor().Find(id)
	assert.NoError(t, err)
	assert.Contains(t, actor.Films, sampleFilm1)
	assert.Contains(t, actor.Films, sampleFilm2)

	actor, err = s.Actor().Find(id + 10)
	assert.Nil(t, actor)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}

func TestActorRepository_GetAll(t *testing.T) {
	s := testdb.New()

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
	s.Actor().Create(actorReq1)
	s.Actor().Create(actorReq2)

	actors, err := s.Actor().GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(actors))
	for _, actor := range actors {
		assert.Contains(t, actor.Films, sampleFilm1)
		assert.Contains(t, actor.Films, sampleFilm2)
	}
}
