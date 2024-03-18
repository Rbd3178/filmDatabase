package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Rbd3178/filmDatabase/internal/app/entities"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
)

// FilmRepository
type FilmRepository struct {
	store *Store
}

// Create
func (r *FilmRepository) Create(f *models.FilmRequest) (id int, done bool, err error) {
	tx, err := r.store.db.Beginx()
	if err != nil {
		return 0, false, errors.Wrap(err, "could not start transaction")
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

	for _, actorID := range f.ActorsIDs {
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

// Delete
func (r *FilmRepository) Delete(id int) (done bool, err error) {
	tx, err := r.store.db.Beginx()
	if err != nil {
		return false, errors.Wrap(err, "could not start transaction")
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

	return r.delete(tx, id)
}

func (r *FilmRepository) delete(tx *sqlx.Tx, id int) (bool, error) {
	_, err := tx.Exec(
		"DELETE FROM films_x_actors WHERE film_id = $1",
		id,
	)
	if err != nil {
		return false, errors.Wrap(err, "delete from films_x_actors")
	}

	res, err := tx.Exec(
		"DELETE FROM films WHERE id = $1",
		id,
	)
	if err != nil {
		return false, errors.Wrap(err, "delete from films")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "rows affected")
	}
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

// Find
func (r *FilmRepository) Find(id int) (actor *models.Film, err error) {
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

	return r.find(tx, id)
}

func (r *FilmRepository) find(tx *sqlx.Tx, id int) (*models.Film, error) {
	var filmInfo entities.Film

	err := tx.Get(
		&filmInfo,
		"SELECT * FROM films WHERE id = $1",
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, errors.Wrap(err, "select film")
	}

	film := models.Film{
		ID:          filmInfo.ID,
		Title:       filmInfo.Title,
		ReleaseDate: filmInfo.ReleaseDate,
		Rating:      filmInfo.Rating,
	}

	err = tx.Select(
		&film.Actors,
		`SELECT
			fxa.actor_id,
			a.name
		FROM 
			films f
		INNER JOIN
			films_x_actors fxa ON fxa.film_id = f.id
		INNER JOIN
			actors a ON a.id = fxa.actor_id
		WHERE f.id = $1`,
		id,
	)

	if err != nil {
		return nil, errors.Wrap(err, "select actors")
	}

	return &film, nil
}

// Modify
func (r *FilmRepository) Modify(id int, f *models.FilmRequest) (done bool, err error) {
	tx, err := r.store.db.Beginx()
	if err != nil {
		return false, errors.Wrap(err, "could not start transaction")
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

	return r.modify(tx, id, f)
}

func (r *FilmRepository) modify(tx *sqlx.Tx, id int, f *models.FilmRequest) (bool, error) {
	if f.Title != "" {
		res, err := tx.Exec(
			"UPDATE films SET title=$1 WHERE id = $2",
			f.Title,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update title")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	if f.Description != "" {
		res, err := tx.Exec(
			"UPDATE films SET description=$1 WHERE id = $2",
			f.Description,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update description")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	if f.ReleaseDate != "" {
		res, err := tx.Exec(
			"UPDATE films SET release_date=$1 WHERE id = $2",
			f.ReleaseDate,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update release date")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	if f.Rating != 0 {
		res, err := tx.Exec(
			"UPDATE films SET rating=$1 WHERE id = $2",
			f.Rating,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update rating")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	_, err := tx.Exec(
		"DELETE FROM films_x_actors WHERE film_id = $1",
		id,
	)
	if err != nil {
		return false, errors.Wrap(err, "delete from films_x_actors")
	}

	for _, actorID := range f.ActorsIDs {
		var exists bool
		err = tx.Get(
			&exists,
			"SELECT EXISTS(SELECT 1 FROM actors WHERE id = $1)",
			actorID,
		)
		if err != nil {
			return false, errors.Wrap(err, "check if actor exists")
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
			return false, errors.Wrap(err, "insert into films_x_actors")
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected films_x_actors")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}

	return true, nil
}
