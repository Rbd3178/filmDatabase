package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

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
		if rowsAffected == 0 {
			return 0, false, errors.Wrap(err, "rows affected films_x_actors")
		}
	}

	return int(id), true, nil
}
