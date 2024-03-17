package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Rbd3178/filmDatabase/internal/app/entities"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
)

// ActorRepository
type ActorRepository struct {
	store *Store
}

// Create
func (r *ActorRepository) Create(a *models.ActorRequest) (id int, err error) {
	tx, err := r.store.db.Beginx()
	if err != nil {
		return 0, errors.Wrap(err, "could not start transaction")
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

	return r.create(tx, a)
}

func (r *ActorRepository) create(tx *sqlx.Tx, a *models.ActorRequest) (int, error) {
	var id int64

	err := tx.Get(
		&id,
		"INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id ",
		a.Name,
		a.Gender,
		a.BirthDate,
	)

	if err != nil {
		return 0, errors.Wrap(err, "insert")
	}

	return int(id), nil
}

// Modify
func (r *ActorRepository) Modify(id int, a *models.ActorRequest) (done bool, err error) {
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

	return r.modify(tx, id, a)
}

func (r *ActorRepository) modify(tx *sqlx.Tx, id int, a *models.ActorRequest) (bool, error) {
	if a.Name != "" {
		res, err := tx.Exec(
			"UPDATE actors SET name=$1 WHERE id = $2",
			a.Name,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update name")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	if a.Gender != "" {
		res, err := tx.Exec(
			"UPDATE actors SET gender=$1 WHERE id = $2",
			a.Gender,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update gender")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	if a.BirthDate != "" {
		res, err := tx.Exec(
			"UPDATE actors SET birth_date=$1 WHERE id = $2",
			a.BirthDate,
			id,
		)
		if err != nil {
			return false, errors.Wrap(err, "update birth date")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return false, errors.Wrap(err, "rows affected")
		}
		if rowsAffected == 0 {
			return false, nil
		}
	}
	
	return true, nil
}

// Delete
func (r *ActorRepository) Delete(id int) (done bool, err error) {
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

func (r *ActorRepository) delete(tx *sqlx.Tx, id int) (bool, error) {
	res, err := tx.Exec(
		"DELETE FROM actors WHERE id = $1",
		id,
	)
	if err != nil {
		return false, errors.Wrap(err, "delete from actors")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "rows affected")
	}
	if rowsAffected == 0 {
		return false, nil
	}

	_, err = tx.Exec(
		"DELETE FROM films_x_actors WHERE actor_id = $1",
		id,
	)
	if err != nil {
		return false, errors.Wrap(err, "delete from films_x_actors")
	}
	return true, nil
}

// Find
func (r *ActorRepository) Find(id int) (actor *models.Actor, err error) {
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

func (r *ActorRepository) find(tx *sqlx.Tx, id int) (*models.Actor, error) {
	var actorInfo entities.Actor

	err := tx.Get(
		&actorInfo,
		"SELECT * FROM actors WHERE id = $1",
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, errors.Wrap(err, "select actor")
	}

	actor := models.Actor{
		ID:        actorInfo.ID,
		Name:      actorInfo.Name,
		Gender:    actorInfo.Gender,
		BirthDate: actorInfo.BirthDate,
	}

	err = tx.Select(
		&actor.Films,
		`SELECT
			fxa.film_id,
			f.title
		FROM 
			actors a
		INNER JOIN
			films_x_actors fxa on fxa.actor_id = a.id
		INNER JOIN
			films f on f.id = fxa.film_id
		WHERE a.id = $1`,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(store.ErrRecordNotFound, "select films")
		}
		return nil, errors.Wrap(err, "select films")
	}

	return &actor, nil
}

// GetAll
func (r *ActorRepository) GetAll() (actor []models.Actor, err error) {
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

	return r.getAll(tx)
}

func (r *ActorRepository) getAll(tx *sqlx.Tx) ([]models.Actor, error) {
	var rawActors = make([]entities.ActorWithFilm, 0)

	err := tx.Select(
		&rawActors,
		`SELECT
			a.id,
			a.name,
			a.gender,
			a.birth_date,
			fxa.film_id,
			f.title
		FROM 
			actors a
		LEFT JOIN
			films_x_actors fxa on fxa.actor_id = a.id
		LEFT JOIN
			films f on f.id = fxa.film_id`,
	)

	if err != nil {
		return nil, errors.Wrap(err, "select")
	}

	actorsMap := make(map[int]models.Actor, 0)
	for _, rawActor := range rawActors {
		actor, ok := actorsMap[rawActor.ID]
		if !ok {
			actorsMap[rawActor.ID] = models.Actor{
				ID: rawActor.ID,
				Name: rawActor.Name,
				Gender: rawActor.Gender,
				BirthDate: rawActor.BirthDate,
			}
		}
		if rawActor.FilmId == nil {
			continue
		}
		actor.Films = append(actor.Films, models.FilmBasic{
			FilmID: *rawActor.FilmId,
			Title: *rawActor.Title,
		})
		actorsMap[rawActor.ID] = actor
	}

	actors := make([]models.Actor, 0, len(actorsMap))
	for _, actor := range actorsMap {
		actors = append(actors, actor)
	}

	return actors, nil
}
