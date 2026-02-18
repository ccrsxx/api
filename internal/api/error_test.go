package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/test"
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

		var res ErrorResponse

		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Error.Message != "An internal server error occurred" {
			t.Errorf("got message %q, want generic server error", res.Error.Message)
		}
	})

	t.Run("PanicError - Development Mode (Stack Trace)", func(t *testing.T) {
		cfg := config.Config()
		originalDev := cfg.IsDevelopment
		cfg.IsDevelopment = true

		defer func() {
			cfg.IsDevelopment = originalDev
		}()

		w := httptest.NewRecorder()

		panicErr := &PanicError{
			Message: "Dev Error",
			Stack:   "goroutine stack trace...",
		}

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

		var res ErrorResponse

		err := json.Unmarshal(w.Body.Bytes(), &res)

		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Error.Message != httpErr.Message {
			t.Errorf("got message %q, want %q", res.Error.Message, httpErr.Message)
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
		w1 := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		HandleHttpError(w1, r, &HttpError{StatusCode: 400, Message: "bad"})

		if w1.Code != 400 {
			t.Errorf("got %d, want 400", w1.Code)
		}

		// Panic error and generic error only returns 500.
		// So we can assert that the status code is 500 to confirm the error handling path was followed.

		w2 := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		HandleHttpError(w2, r, &PanicError{Message: "crash"})

		if w2.Code != 500 {
			t.Errorf("got %d, want 500", w2.Code)
		}

		w3 := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		HandleHttpError(w3, r, errors.New("generic"))

		if w3.Code != 500 {
			t.Errorf("got %d, want 500", w3.Code)
		}
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
