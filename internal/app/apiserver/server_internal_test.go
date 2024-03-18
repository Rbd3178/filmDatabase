package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rbd3178/filmDatabase/internal/app/store/testdb"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleRegister(t *testing.T) {
	s := newServer(testdb.New())
	var tests = []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"login":    "somebody",
				"password": "secret",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"login":    "nottoolong",
				"password": "short",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "login already taken",
			payload: map[string]interface{}{
				"login":    "somebody",
				"password": "qwerty",
			},
			expectedCode: http.StatusConflict,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/register", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

// Primarily a way to test authorization, for other hadlers only check 403 response
func TestServer_HandleUsers(t *testing.T) {
	s := newServer(testdb.New())

	var tests = []struct {
		name         string
		login        string
		password     string
		expectedCode int
	}{
		{
			name:         "Unregistered",
			login:        "nobody",
			password:     "pass",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Registered, incorrect password",
			login:        "normal",
			password:     "incorrect",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Valid auth info, not admin",
			login:        "normal",
			password:     "correct",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Admin",
			login:        "admin",
			password:     "adminpass",
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			req.SetBasicAuth(tc.login, tc.password)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotNil(t, rec.Body)
		})
	}
}

func TestServer_HandleActorsGet(t *testing.T) {
	s := newServer(testdb.New())

	var tests = []struct {
		name         string
		method       string
		login        string
		password     string
		expectedCode int
	}{
		{
			name:         "Request with incorrect method",
			method:       http.MethodPut,
			login:        "normal",
			password:     "correct",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Request by normal user",
			method:       http.MethodGet,
			login:        "normal",
			password:     "correct",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Request by admin",
			method:       http.MethodGet,
			login:        "admin",
			password:     "adminpass",
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "/actors", nil)
			req.SetBasicAuth(tc.login, tc.password)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotNil(t, rec.Body)
		})
	}
}

func TestServer_HandleActorsPost(t *testing.T) {
	s := newServer(testdb.New())

	var tests = []struct {
		name         string
		method       string
		login        string
		password     string
		payload      any
		expectedCode int
	}{
		{
			name:         "Request with incorrect method",
			method:       http.MethodPut,
			login:        "normal",
			password:     "correct",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "Request by normal user",
			method:   http.MethodPost,
			login:    "normal",
			password: "correct",
			payload: map[string]interface{}{
				"name":       "Name",
				"gender":     "Gender",
				"birth_date": "2024-03-18",
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:     "Request by admin",
			method:   http.MethodPost,
			login:    "admin",
			password: "adminpass",
			payload: map[string]interface{}{
				"name":       "Name",
				"gender":     "Gender",
				"birth_date": "2024-03-18",
			},
			expectedCode: http.StatusCreated,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(tc.method, "/actors", b)
			req.SetBasicAuth(tc.login, tc.password)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotNil(t, rec.Body)
		})
	}
}
