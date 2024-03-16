package postgres_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "user=postgres password=postgres host=localhost dbname=filmdb_test sslmode=disable"
	}

	os.Exit(m.Run())
}
