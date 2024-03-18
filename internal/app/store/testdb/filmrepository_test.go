package testdb_test

import (
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/Rbd3178/filmDatabase/internal/app/store/testdb"
	"github.com/stretchr/testify/assert"
)

func TestFilmRepository_Create(t *testing.T) {
	s := testdb.New()

	filmReq := &models.FilmRequest{
		Title:       "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      6.8,
		ActorsIDs:   []int{1, 2, 5},
	}

	id, err := s.Film().Create(filmReq)
	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestFilmRepository_Delete(t *testing.T) {
	s := testdb.New()

	filmReq := &models.FilmRequest{
		Title:       "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      7.8,
		ActorsIDs:   []int{2, 3},
	}
	filmId, _ := s.Film().Create(filmReq)

	done, err := s.Film().Delete(filmId)
	assert.NoError(t, err)
	assert.True(t, done)

	done, err = s.Film().Delete(filmId)
	assert.NoError(t, err)
	assert.False(t, done)
}

func TestFilmRepository_Modify(t *testing.T) {
	s := testdb.New()

	filmReq := &models.FilmRequest{
		Title:       "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      6.8,
		ActorsIDs:   []int{2},
	}
	filmID, _ := s.Film().Create(filmReq)

	filmReqMod := &models.FilmRequest{
		Description: "Even more detailed description",
		Rating:      6.9,
		ActorsIDs:   []int{1, 3},
	}

	done, err := s.Film().Modify(filmID, filmReqMod)
	assert.NoError(t, err)
	assert.True(t, done)

	done, err = s.Film().Modify(filmID+10, filmReqMod)
	assert.NoError(t, err)
	assert.False(t, done)
}

func TestFilmRepository_Find(t *testing.T) {
	s := testdb.New()

	filmReq := &models.FilmRequest{
		Title:       "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      6.8,
		ActorsIDs:   []int{2, 3},
	}
	filmID, _ := s.Film().Create(filmReq)

	film, err := s.Film().Find(filmID)
	assert.NoError(t, err)
	assert.Contains(t, film.Actors, models.ActorBasic{
		ActorID: 2,
		Name: "Second Actor",
	})
	assert.Contains(t, film.Actors, models.ActorBasic{
		ActorID: 3,
		Name: "Third Actor",
	})

	film, err = s.Film().Find(filmID + 10)
	assert.Nil(t, film)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}

func TestFilmRepository_GetAll(t *testing.T) {
	s := testdb.New()

	filmReq1 := &models.FilmRequest{
		Title:       "Cool title",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      7.8,
		ActorsIDs:   []int{1, 2, 3},
	}
	filmReq2 := &models.FilmRequest{
		Title:       "Cool title 2",
		Description: "Detailed description",
		ReleaseDate: "2020-01-01",
		Rating:      6.8,
		ActorsIDs:   []int{},
	}
	filmID1, _ := s.Film().Create(filmReq1)
	filmID2, _ := s.Film().Create(filmReq2)

	films, err := s.Film().GetAll("rating", "desc", "cool", "third")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(films))
	for _, film := range films {
		if film.ID == filmID2 {
			assert.Contains(t, film.Actors, models.ActorBasic{ActorID: 1, Name: "First Actor"})
			assert.Contains(t, film.Actors, models.ActorBasic{ActorID: 2, Name: "Second Actor"})
			assert.Contains(t, film.Actors, models.ActorBasic{ActorID: 3, Name: "Third Actor"})
			assert.Equal(t, filmReq2.Title, film.Title)
		}
		if film.ID == filmID1 {
			assert.Empty(t, film.Actors)
		}
	}
}