package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Rbd3178/filmDatabase/internal/app/hasher"
	"github.com/Rbd3178/filmDatabase/internal/app/models"
	"github.com/Rbd3178/filmDatabase/internal/app/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router   *http.ServeMux
	logger   *logrus.Logger
	database store.Store
}

func newServer(database store.Store) *server {
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
	s.router.HandleFunc("/register", s.handleRegister)
	s.router.HandleFunc("/users", s.handleUsers)
	s.router.HandleFunc("/actors", s.handleActors)
	s.router.HandleFunc("/actors/", s.handleActorsID)
	s.router.HandleFunc("/films", s.handleFilms)
	s.router.HandleFunc("/films/", s.handleFilmsID)
}

func (s *server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.registerUser(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) registerUser(w http.ResponseWriter, r *http.Request) {
	req := &models.UserRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !req.Validate() {
		http.Error(w, "Invalid login or password", http.StatusUnprocessableEntity)
		return
	}
	done, err := s.database.User().Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !done {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User successfully registered"))
}

func (s *server) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	registered, isAdmin := s.authenticateUser(w, r)
	if !registered {
		return
	}
	if !isAdmin {
		http.Error(w, "Not enough rights", http.StatusForbidden)
		return
	}

	s.getUsers(w)
}

func (s *server) authenticateUser(w http.ResponseWriter, r *http.Request) (bool, bool) {
	login, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		http.Error(w, "No basic auth present", http.StatusUnauthorized)
		return false, false
	}

	user, err := s.database.User().Find(login)
	if err != nil {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		http.Error(w, "Incorrect login", http.StatusUnauthorized)
		return false, false
	}

	if !hasher.CheckPasswordHash(password, user.HashedPassword) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return false, false
	}

	return true, user.IsAdmin
}

func (s *server) getUsers(w http.ResponseWriter) {
	users, err := s.database.User().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) handleActors(w http.ResponseWriter, r *http.Request) {
	registered, isAdmin := s.authenticateUser(w, r)
	if !registered {
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getActors(w)

	case http.MethodPost:
		if !isAdmin {
			http.Error(w, "Not enough rights", http.StatusForbidden)
			return
		}
		s.addActor(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) getActors(w http.ResponseWriter) {
	actors, err := s.database.Actor().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) addActor(w http.ResponseWriter, r *http.Request) {
	req := &models.ActorRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !req.ValidateForInsert() {
		http.Error(w, "Invalid fields in payload", http.StatusUnprocessableEntity)
		return
	}
	id, err := s.database.Actor().Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/actors/%d", id))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Actor successfully added"))
}

func (s *server) handleActorsID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	registered, isAdmin := s.authenticateUser(w, r)
	if !registered {
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.findActor(w, r, id)

	case http.MethodPatch:
		if !isAdmin {
			http.Error(w, "Not enough rights", http.StatusForbidden)
			return
		}
		s.modifyActor(w, r, id)

	case http.MethodDelete:
		if !isAdmin {
			http.Error(w, "Not enough rights", http.StatusForbidden)
			return
		}
		s.deleteActor(w, r, id)
	}
}

func (s *server) findActor(w http.ResponseWriter, r *http.Request, id int) {
	actor, err := s.database.Actor().Find(id)
	if err == store.ErrRecordNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) modifyActor(w http.ResponseWriter, r *http.Request, id int) {
	req := &models.ActorRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !req.ValidateForUpdate() {
		http.Error(w, "Invalid fields in payload", http.StatusUnprocessableEntity)
		return
	}

	done, err := s.database.Actor().Modify(id, req)
	if !done {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Actor information successfully modified"))
}

func (s *server) deleteActor(w http.ResponseWriter, r *http.Request, id int) {
	done, err := s.database.Actor().Delete(id)
	if !done {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handleFilms(w http.ResponseWriter, r *http.Request) {
	registered, isAdmin := s.authenticateUser(w, r)
	if !registered {
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getFilms(w, r)

	case http.MethodPost:
		if !isAdmin {
			http.Error(w, "Not enough rights", http.StatusForbidden)
			return
		}
		s.addFilm(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) addFilm(w http.ResponseWriter, r *http.Request) {
	req := &models.FilmRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !req.ValidateForInsert() {
		http.Error(w, "Invalid fields in payload", http.StatusUnprocessableEntity)
		return
	}

	id, err := s.database.Film().Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/films/%d", id))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Film successfully added"))
}

func (s *server) getFilms(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	orderBy := query.Get("orderby")
	if orderBy == "" {
		orderBy = "rating"
	}

	order := query.Get("order")
	if order == "" {
		order = "desc"
	}

	if orderBy != "rating" && orderBy != "title" && orderBy != "release_date" || order != "asc" && order != "desc" {
		http.Error(w, "invalid query parameters", http.StatusBadRequest)
		return
	}

	searchTitle := query.Get("searchtitle")
	searchActor := query.Get("searchactor")

	films, err := s.database.Film().GetAll(orderBy, order, searchTitle, searchActor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(films)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) handleFilmsID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	registered, isAdmin := s.authenticateUser(w, r)
	if !registered {
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.findFilm(w, r, id)

	case http.MethodPatch:
		if !isAdmin {
			http.Error(w, "Not enough rights", http.StatusForbidden)
			return
		}
		s.modifyFilm(w, r, id)

	case http.MethodDelete:
		if !isAdmin {
			http.Error(w, "Not enough rights", http.StatusForbidden)
			return
		}
		s.deleteFilm(w, r, id)
	}
}

func (s *server) findFilm(w http.ResponseWriter, r *http.Request, id int) {
	actor, err := s.database.Film().Find(id)
	if err == store.ErrRecordNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) modifyFilm(w http.ResponseWriter, r *http.Request, id int) {
	req := &models.FilmRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !req.ValidateForUpdate() {
		http.Error(w, "Invalid fields in payload", http.StatusUnprocessableEntity)
		return
	}

	done, err := s.database.Film().Modify(id, req)
	if !done {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Film information successfully modified"))
}

func (s *server) deleteFilm(w http.ResponseWriter, r *http.Request, id int) {
	done, err := s.database.Film().Delete(id)
	if !done {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}