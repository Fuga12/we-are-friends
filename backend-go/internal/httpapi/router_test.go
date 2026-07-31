package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fuga12/we-are-friends/backend-go/internal/config"
)

func TestHealthz(t *testing.T) {
	r := NewRouter(config.Config{CORSOrigins: []string{"http://localhost:5173"}}, nil, testMedia{})
	req := httptest.NewRequest(http.MethodGet, "/api/healthz", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if body := res.Body.String(); body != "{\"ok\":true}\n" {
		t.Fatalf("body = %q, want health response", body)
	}
}

type testMedia struct{}

func (testMedia) URL(_ context.Context, key string) (string, error) {
	return "https://media.test/" + key, nil
}
