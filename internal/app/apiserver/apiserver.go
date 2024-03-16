package apiserver

import (
	"net/http"

	"github.com/Rbd3178/filmDatabase/internal/app/store/postgres"
	"github.com/jmoiron/sqlx"
)

// Start
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	database := postgres.New(db)
	srv := newServer(database)
	return http.ListenAndServe(config.Port, srv)
}

func newDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
