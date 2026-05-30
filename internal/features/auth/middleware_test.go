package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
)

func TestMiddleware_IsAuthorizedFromBearer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "middleware-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer middleware-secret")

		mw.IsAuthorizedFromBearer(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "middleware-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer wrong-key")

		mw.IsAuthorizedFromBearer(handler).ServeHTTP(w, r)

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
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "query-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=query-secret", nil)

		mw.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "query-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)

		mw.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}

func TestMiddleware_IsAuthorizedFromBearerOrQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success via Bearer Header (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "combo-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer combo-secret")

		mw.IsAuthorizedFromBearerOrQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Success via Query Token (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "combo-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=combo-secret", nil)

		mw.IsAuthorizedFromBearerOrQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "combo-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)
		r.Header.Set("Authorization", "Bearer wrong-secret")

		mw.IsAuthorizedFromBearerOrQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}
