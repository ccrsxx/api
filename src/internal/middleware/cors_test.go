package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
)

func TestCors(t *testing.T) {
	// Setup config for tests
	// We modify the global config instance directly since it returns a pointer
	cfg := config.Env()

	originalOrigins := cfg.AllowedOrigins
	cfg.AllowedOrigins = []string{"https://allowed.com", "http://localhost:3000"}

	defer func() { cfg.AllowedOrigins = originalOrigins }()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Allowed Origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.Header.Set("Origin", "https://allowed.com")

		w := httptest.NewRecorder()

		Cors(nextHandler).ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "https://allowed.com" {
			t.Errorf("expected Access-Control-Allow-Origin header")
		}

		if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
			t.Errorf("expected Access-Control-Allow-Credentials header")
		}
	})

	t.Run("Disallowed Origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.Header.Set("Origin", "https://hacker.com")

		w := httptest.NewRecorder()

		Cors(nextHandler).ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("expected no Access-Control-Allow-Origin header for disallowed origin")
		}
	})

	t.Run("Preflight OPTIONS Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)

		req.Header.Set("Origin", "https://allowed.com")

		w := httptest.NewRecorder()

		Cors(nextHandler).ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("got status %d, want %d for OPTIONS", w.Code, http.StatusNoContent)
		}

		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("expected Access-Control-Allow-Methods header")
		}
	})
}
