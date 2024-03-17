package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	switch r.Method {
	case http.MethodGet:
		registered, _ := s.authenticateUser(w, r)
		if !registered {
			return
		}
		s.getActors(w)

	case http.MethodPost:
		registered, isAdmin := s.authenticateUser(w, r)
		if !registered {
			return
		}
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
	if !req.Validate() {
		http.Error(w, "Invalid fields in payload", http.StatusUnprocessableEntity)
		return
	}
	id, err := s.database.Actor().Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location: ", fmt.Sprintf("/actors/%d", id))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Actor successfully added"))
}
