package api

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// faultyWriter mocks a broken network connection (always fails Write).
// Since it is in package 'api', it is available to error_test.go as well.
type faultyWriter struct {
	*httptest.ResponseRecorder
}

func (fw *faultyWriter) Write(b []byte) (int, error) {
	return 0, errors.New("forced network error")
}

type testUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNewSuccessResponse(t *testing.T) {
	t.Run("Wraps Struct in Data Field", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := testUser{Name: "John", Age: 30}

		err := NewSuccessResponse(w, http.StatusOK, data)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}

		expected := `{"data":{"name":"John","age":30}}`

		if strings.TrimSpace(w.Body.String()) != expected {
			t.Errorf("got body %q, want %q", w.Body.String(), expected)
		}
	})

	t.Run("Wraps Pointer to Struct in Data Field", func(t *testing.T) {
		w := httptest.NewRecorder()

		data := &testUser{Name: "Jane", Age: 25}

		err := NewSuccessResponse(w, http.StatusCreated, data)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := `{"data":{"name":"Jane","age":25}}`

		if strings.TrimSpace(w.Body.String()) != expected {
			t.Errorf("got body %q, want %q", w.Body.String(), expected)
		}
	})

	t.Run("Does Not Wrap Maps", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"foo": "bar"}

		// Maps are treated as direct responses
		err := NewSuccessResponse(w, http.StatusOK, data)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := `{"foo":"bar"}`

		if strings.TrimSpace(w.Body.String()) != expected {
			t.Errorf("got body %q, want %q", w.Body.String(), expected)
		}
	})

	t.Run("Does Not Wrap Slices", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []string{"a", "b"}

		err := NewSuccessResponse(w, http.StatusOK, data)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := `["a","b"]`

		if strings.TrimSpace(w.Body.String()) != expected {
			t.Errorf("got body %q, want %q", w.Body.String(), expected)
		}
	})
}

func TestNewErrorResponse(t *testing.T) {
	t.Run("Standard Error Response", func(t *testing.T) {
		w := httptest.NewRecorder()
		msg := "Invalid input"
		details := []string{"field x is required"}
		id := "test-id"

		err := NewErrorResponse(w, http.StatusBadRequest, msg, details, id)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if w.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
		}

		var response ErrorResponse

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Message != msg {
			t.Errorf("got message %q, want %q", response.Error.Message, msg)
		}

		if response.Error.ID != id {
			t.Errorf("got id %q, want %q", response.Error.ID, id)
		}

		if len(response.Error.Details) != 1 || response.Error.Details[0] != details[0] {
			t.Errorf("got details %v, want %v", response.Error.Details, details)
		}
	})

	t.Run("Handles Nil Details", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Passing nil for details should result in empty slice []
		_ = NewErrorResponse(w, http.StatusNotFound, "Not Found", nil, "id")

		var response ErrorResponse
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		if response.Error.Details == nil {
			t.Error("expected empty slice for details, got nil")
		}

		if len(response.Error.Details) != 0 {
			t.Error("expected empty slice length 0")
		}
	})
}

// TestInternalNewResponse tests the private 'newResponse' function directly
// to cover edge cases like network failures.
func TestInternalNewResponse(t *testing.T) {
	t.Run("Write Failure (Network Error)", func(t *testing.T) {
		w := &faultyWriter{ResponseRecorder: httptest.NewRecorder()}
		data := map[string]string{"foo": "bar"}

		err := newResponse(w, http.StatusOK, data)

		if err == nil {
			t.Fatal("expected error from Write failure, got nil")
		}

		if !strings.Contains(err.Error(), "write response error") {
			t.Errorf("expected 'write response error', got %v", err)
		}
	})

	t.Run("Marshal Failure (Fallback Success)", func(t *testing.T) {
		w := httptest.NewRecorder()

		data := map[string]float64{"val": math.Inf(1)} // Infinity fails Marshal

		err := newResponse(w, http.StatusOK, data)

		if err == nil {
			t.Fatal("expected error from marshal failure, got nil")
		}

		// It should contain the marshal error...
		if !strings.Contains(err.Error(), "marshal response error") {
			t.Errorf("expected 'marshal response error' in %v", err)
		}

		// ...but NOT the fallback error (because fallback succeeded)
		if strings.Contains(err.Error(), "marshal fallback error") {
			t.Errorf("unexpected 'marshal fallback error' in %v", err)
		}

		// Verify it wrote the fallback 500 status
		if w.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want 500", w.Code)
		}
	})

	t.Run("Double Failure (Marshal + Write)", func(t *testing.T) {
		w := &faultyWriter{ResponseRecorder: httptest.NewRecorder()}

		// Infinity fails json.Marshal, triggering the fallback logic
		data := map[string]float64{"val": math.Inf(1)}

		err := newResponse(w, http.StatusOK, data)

		if err == nil {
			t.Fatal("expected error from double failure, got nil")
		}

		// The error should contain both the marshal error AND the fallback write error
		if !strings.Contains(err.Error(), "marshal response error") {
			t.Errorf("expected 'marshal response error' in %v", err)
		}

		if !strings.Contains(err.Error(), "marshal fallback error") {
			t.Errorf("expected 'marshal fallback error' in %v", err)
		}
	})
}
