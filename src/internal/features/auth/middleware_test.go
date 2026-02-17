package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
)

func TestMiddleware_IsAuthorized(t *testing.T) {

	originalKey := config.Env().SecretKey

	defer func() {
		config.Env().SecretKey = originalKey
	}()

	config.Env().SecretKey = "middleware-secret"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer middleware-secret") // Matches mocked env

		w := httptest.NewRecorder()

		Middleware.IsAuthorized(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want status 200, got %d", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		
		r.Header.Set("Authorization", "Bearer wrong-key") // Does NOT match mocked env

		w := httptest.NewRecorder()

		Middleware.IsAuthorized(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("want status 401, got %d", w.Code)
		}
	})
}

func TestMiddleware_IsAuthorizedFromQuery(t *testing.T) {
	originalKey := config.Env().SecretKey

	defer func() {
		config.Env().SecretKey = originalKey
	}()

	config.Env().SecretKey = "query-secret"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/?token=query-secret", nil)
		w := httptest.NewRecorder()

		Middleware.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want status 200, got %d", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)
		w := httptest.NewRecorder()

		Middleware.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("want status 401, got %d", w.Code)
		}
	})
}
