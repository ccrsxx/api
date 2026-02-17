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
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer middleware-secret")

		Middleware.IsAuthorized(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.Header.Set("Authorization", "Bearer wrong-key")

		Middleware.IsAuthorized(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
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
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=query-secret", nil)

		Middleware.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)

		Middleware.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}
