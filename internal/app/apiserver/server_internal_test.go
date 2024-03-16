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
		name string
		payload any
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
			name: "invalid payload",
			payload: "invalid",
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