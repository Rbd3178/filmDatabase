package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Rbd3178/filmDatabase/internal/app/entities"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
)

// FilmRepository
type FilmRepository struct {
	store *Store
}

// Create
func (r *FilmRepository) Create(f *models.FilmRequest) (id int, done bool, err error) {
	tx, err := r.store.db.Beginx()
	if err != nil {
		return 0, false, errors.Wrap(err, "can not start transaction")
	}

	defer func() {
		if err != nil {
			errRb := tx.Rollback()
			if errRb != nil {
				errors.Wrap(err, "error during rollback")
				return
			}

			return
		}
		err = tx.Commit()
	}()

	return r.create(tx, f)
}

func (r *FilmRepository) create(tx *sqlx.Tx, f *models.FilmRequest) (int, bool, error) {
	var id int64
	err := tx.Get(
		&id,
		"INSERT INTO films (title, description, release_date, rating) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id ",
		f.Title,
		f.Description,
		f.ReleaseDate,
		f.Rating,
	)

	if err != nil {
		return 0, false, errors.Wrap(err, "insert into films")
	}

	for _, actorID := range f.Actors_IDs {
		var exists bool
		err = tx.Get(
			&exists,
			"SELECT EXISTS(SELECT 1 FROM actors WHERE id = $1)",
			actorID,
		)
		if err != nil {
			return 0, false, errors.Wrap(err, "check if actor exists")
		}
		if !exists {
			continue
		}

		res, err := tx.Exec(
			"INSERT INTO films_x_actors (film_id, actor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			id,
			actorID,
		)

		if err != nil {
			return 0, false, errors.Wrap(err, "insert into films_x_actors")
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, false, errors.Wrap(err, "rows affected films_x_actors")
		}
		if rowsAffected == 0 {
			return 0, false, nil
		}
	}

	return int(id), true, nil
}

// GetAll
func (r *FilmRepository) GetAll(orderBy string, order string, searchTitle string, searchActor string) (films []models.Film, err error) {
	tx, err := r.store.db.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "could not start transaction")
	}

	defer func() {
		if err != nil {
			errRb := tx.Rollback()
			if errRb != nil {
				err = errors.Wrap(err, "error during rollback")
				return
			}

			return
		}

		err = tx.Commit()
	}()

	return r.getAll(tx, orderBy, order, searchTitle, searchActor)
}

func (r *FilmRepository) getAll(tx *sqlx.Tx, orderBy string, order string, searchTitle string, searchActor string) ([]models.Film, error) {
	var rawFilms = make([]entities.FilmWithActor, 0)
	query := `SELECT
				f.id,
				f.title,
				f.description,
				f.release_date,
				f.rating,
				fxa.actor_id,
				a.name 
			FROM 
				films f
			LEFT JOIN
				films_x_actors fxa ON fxa.film_id = f.id
			LEFT JOIN
				actors a ON a.id = fxa.actor_id` + fmt.Sprintf(" WHERE f.title ILIKE '%%%s%%' ORDER BY f.%s %s, f.id ASC", searchTitle, orderBy, order)

	err := tx.Select(&rawFilms, query)

	if err != nil {
		return nil, errors.Wrap(err, "select")
	}

	var idsActorSearch = make([]int, 0)
	var idsMap = make(map[int]struct{})
	if searchActor != "" {
		err = tx.Select(
			&idsActorSearch,
			`SELECT
			f.id
		FROM
			films f
		LEFT JOIN
			films_x_actors fxa ON fxa.film_id = f.id
		LEFT JOIN
			actors a ON a.id = fxa.actor_id`+fmt.Sprintf(" WHERE a.name ILIKE '%%%s%%'", searchActor),
		)

		if err != nil {
			return nil, errors.Wrap(err, "select ids actor search")
		}

		for _, id := range idsActorSearch {
			idsMap[id] = struct{}{}
		}
	}
	films := make([]models.Film, 0)
	curID := -1
	for _, rawFilm := range rawFilms {
		if _, ok := idsMap[rawFilm.ID]; !ok && searchActor != "" {
			continue
		}
		if rawFilm.ID != curID {
			films = append(films, models.Film{
				ID:          rawFilm.ID,
				Title:       rawFilm.Title,
				Description: rawFilm.Description,
				ReleaseDate: rawFilm.ReleaseDate,
				Rating:      rawFilm.Rating,
			})
			curID = rawFilm.ID
		}
		if rawFilm.ActorID == nil {
			continue
		}
		films[len(films)-1].Actors = append(films[len(films)-1].Actors, models.ActorBasic{
			ActorID: *rawFilm.ActorID,
			Name:    *rawFilm.Name,
		})
	}

	return films, nil
}
