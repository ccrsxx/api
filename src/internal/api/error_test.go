package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
)

func TestHandleHttpError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	t.Run("PanicError - Standard (Production)", func(t *testing.T) {
		w := httptest.NewRecorder()

		panicErr := &PanicError{
			Message: "Something crashed",
			Stack:   "trace...",
			Value:   "nil pointer",
		}

		HandleHttpError(w, r, panicErr)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
		}

		var resp ErrorResponse

		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		// In production, message is generic
		if resp.Error.Message != "An internal server error occurred" {
			t.Errorf("got message %q, want generic server error", resp.Error.Message)
		}
	})

	t.Run("PanicError - Development Mode (Stack Trace)", func(t *testing.T) {
		cfg := config.Config()
		originalDev := cfg.IsDevelopment
		cfg.IsDevelopment = true

		defer func() { cfg.IsDevelopment = originalDev }() // Reset after test

		w := httptest.NewRecorder()

		panicErr := &PanicError{
			Message: "Dev Error",
			Stack:   "goroutine stack trace...",
		}

		// This triggers the fmt.Printf path in your code
		HandleHttpError(w, r, panicErr)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("HttpError - Custom Status", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpErr := &HttpError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ID",
			Details:    []string{"id must be uuid"},
		}

		HandleHttpError(w, r, httpErr)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", w.Code)
		}

		var resp ErrorResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Error.Message != httpErr.Message {
			t.Errorf("got message %q, want %q", resp.Error.Message, httpErr.Message)
		}
	})

	t.Run("Generic Error - Defaults to 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		genericErr := errors.New("database connection failed")

		HandleHttpError(w, r, genericErr)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want 500", w.Code)
		}
	})

	t.Run("Write Failures (Triggers logErrorResponse)", func(t *testing.T) {
		w1 := &faultyWriter{ResponseRecorder: httptest.NewRecorder()}

		HandleHttpError(w1, r, &PanicError{Message: "crash"})

		w2 := &faultyWriter{ResponseRecorder: httptest.NewRecorder()}

		HandleHttpError(w2, r, &HttpError{StatusCode: 400, Message: "bad"})

		w3 := &faultyWriter{ResponseRecorder: httptest.NewRecorder()}

		HandleHttpError(w3, r, errors.New("generic"))
	})
}

func TestErrorTypes(t *testing.T) {
	t.Run("PanicError.Error()", func(t *testing.T) {
		pe := &PanicError{Message: "panic msg"}

		if pe.Error() != "panic msg" {
			t.Errorf("got %q, want %q", pe.Error(), "panic msg")
		}
	})

	t.Run("HttpError.Error()", func(t *testing.T) {
		he := &HttpError{Message: "http msg"}

		if he.Error() != "http msg" {
			t.Errorf("got %q, want %q", he.Error(), "http msg")
		}
	})
}
