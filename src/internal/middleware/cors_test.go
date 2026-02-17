package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
)

func TestCors(t *testing.T) {
	cfg := config.Env()

	originalOrigins := cfg.AllowedOrigins

	defer func() {
		cfg.AllowedOrigins = originalOrigins
	}()

	cfg.AllowedOrigins = []string{"https://allowed.com", "http://localhost:3000"}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Allowed Origin", func(t *testing.T) {
		tests := []struct {
			name    string
			origin  string
			allowed bool
		}{
			{
				name:    "Allowed Origin 1",
				origin:  "https://allowed.com",
				allowed: true,
			},
			{
				name:    "Allowed Origin 2",
				origin:  "http://localhost:3000",
				allowed: true,
			},
			{
				name:    "Disallowed Origin",
				origin:  "https://hacker.com",
				allowed: false,
			},
			{
				name:    "No Origin Header",
				origin:  "",
				allowed: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				r.Header.Set("Origin", tt.origin)

				w := httptest.NewRecorder()

				Cors(nextHandler).ServeHTTP(w, r)

				if tt.allowed {
					if w.Header().Get("Access-Control-Allow-Origin") != tt.origin {
						t.Errorf("want Access-Control-Allow-Origin header to be %s", tt.origin)
					}

					if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
						t.Errorf("want Access-Control-Allow-Credentials header")
					}
				} else {
					if w.Header().Get("Access-Control-Allow-Origin") != "" {
						t.Errorf("want no Access-Control-Allow-Origin header for disallowed origin")
					}

					if w.Header().Get("Access-Control-Allow-Credentials") != "" {
						t.Errorf("want no Access-Control-Allow-Credentials header for disallowed origin")
					}
				}
			})
		}
	})

	t.Run("Preflight OPTIONS Request", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "/", nil)

		r.Header.Set("Origin", "https://allowed.com")

		w := httptest.NewRecorder()

		Cors(nextHandler).ServeHTTP(w, r)

		if w.Code != http.StatusNoContent {
			t.Errorf("got status %d, want %d for OPTIONS", w.Code, http.StatusNoContent)
		}

		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("want Access-Control-Allow-Methods header")
		}
	})
}
