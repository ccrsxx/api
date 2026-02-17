package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)

		if _, err := w.Write([]byte("I am a teapot")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	r := httptest.NewRequest(http.MethodGet, "/tea", nil)

	w := httptest.NewRecorder()

	Logging(nextHandler).ServeHTTP(w, r)

	if w.Code != http.StatusTeapot {
		t.Errorf("got status %d, want 418", w.Code)
	}

	if w.Body.String() != "I am a teapot" {
		t.Errorf("got body %q, want 'I am a teapot'", w.Body.String())
	}
}

func TestWrappedWriter_Unwrap(t *testing.T) {
	// This covers the Unwrap method in logging.go
	w := httptest.NewRecorder()

	wrapped := &wrappedWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Verify that unwrapping returns the original ResponseWriter
	if wrapped.Unwrap() != w {
		t.Error("Unwrap() did not return the underlying ResponseWriter")
	}
}
