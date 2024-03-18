package postgres

import (
	"database/sql"

	"github.com/Rbd3178/filmDatabase/internal/app/hasher"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// UserRepository
type UserRepository struct {
	store *Store
}

// Create
func (r *UserRepository) Create(u *models.UserRequest) (bool, error) {
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

	return r.create(tx, u)
}

func (r *UserRepository) create(tx *sqlx.Tx, u *models.UserRequest) (bool, error) {
	hashedPassword, err := hasher.HashPassword(u.Password)
	if err != nil {
		return false, errors.Wrap(err, "encryption")
	}

	res, err := tx.Exec(
		"INSERT INTO users (login, hashed_password, is_admin) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		u.Login,
		hashedPassword,
		false,
	)
	if err != nil {
		return false, errors.Wrap(err, "insert")
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
func (r *UserRepository) Find(login string) (*models.User, error) {
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

	return r.find(tx, login)
}

func (r *UserRepository) find(tx *sqlx.Tx, login string) (*models.User, error) {
	u := &models.User{}
	err := tx.Get(
		u,
		"SELECT * FROM users WHERE login = $1",
		login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// GetAll
func (r *UserRepository) GetAll() ([]models.User, error) {
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

func (r *UserRepository) getAll(tx *sqlx.Tx) ([]models.User, error) {
	users := make([]models.User, 0)
	err := tx.Select(
		&users,
		"SELECT * FROM users",
	)
	if err != nil {
		return nil, errors.Wrap(err, "select")
	}
	return users, nil
}
