package apiserver

import (
	"net/http"

	"github.com/Rbd3178/filmDatabase/internal/app/dbs"
	"github.com/sirupsen/logrus"
)

type server struct {
	router   *http.ServeMux
	logger   *logrus.Logger
	database dbs.DBS
}

func newServer(database dbs.DBS) *server {
	s := &server{
		router:   http.NewServeMux(),
		logger:   logrus.New(),
		database: database,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsers())
}

func (s *server) handleUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}
