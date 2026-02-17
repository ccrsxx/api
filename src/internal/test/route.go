package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RouteTestCase struct {
	Path       string
	Host       string
	Method     string
	Headers    http.Header
	StatusCode int // Optional: want status code (defaults to checking for non-404)
}

func AssertRoutes(t *testing.T, mux *http.ServeMux, tests []RouteTestCase) {
	t.Helper()

	for _, tt := range tests {
		testName := fmt.Sprintf("Test route %s with method %s", tt.Method, tt.Path)

		t.Run(testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.Method, tt.Path, nil)

			if tt.Host != "" {
				r.Host = tt.Host
			}

			for key, values := range tt.Headers {
				for _, value := range values {
					r.Header.Add(key, value)
				}
			}

			mux.ServeHTTP(w, r)

			// If StatusCode is not set, just check that the route is registered (not 404)
			if tt.StatusCode == 0 {
				if w.Code == http.StatusNotFound {
					t.Errorf("got 404 for route %s %s, want it to be registered", tt.Method, tt.Path)
				}

				return
			}

			if w.Code != tt.StatusCode {
				t.Errorf("got status %d for route %s %s, want %d", w.Code, tt.Method, tt.Path, tt.StatusCode)
			}
		})
	}
}
