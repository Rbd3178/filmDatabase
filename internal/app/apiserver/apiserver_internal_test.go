package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_HandleHello(t *testing.T) {
	s := New(NewConfig())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	s.handleHello().ServeHTTP(rec, req)
	want := "Hello"
	got := rec.Body.String()
	if got != want {
		t.Errorf("Wanted \"%s\", got \"%s\"", want, got)
	}
}
