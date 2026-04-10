package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_IsAuthorized(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		mw := NewMiddleware(NewService(ServiceConfig{SecretKey: "middleware-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer middleware-secret")

		mw.IsAuthorized(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := NewMiddleware(NewService(ServiceConfig{SecretKey: "middleware-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer wrong-key")

		mw.IsAuthorized(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}

func TestMiddleware_IsAuthorizedFromQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		mw := NewMiddleware(NewService(ServiceConfig{SecretKey: "query-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=query-secret", nil)

		mw.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := NewMiddleware(NewService(ServiceConfig{SecretKey: "query-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)

		mw.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}
