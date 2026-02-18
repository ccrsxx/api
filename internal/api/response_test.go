package api

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

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
			t.Fatalf("unwanted error: %v", err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
		}

		want := `{"data":{"name":"John","age":30}}`

		if strings.TrimSpace(w.Body.String()) != want {
			t.Errorf("got body %q, want %q", w.Body.String(), want)
		}
	})

	t.Run("Wraps Pointer to Struct in Data Field", func(t *testing.T) {
		w := httptest.NewRecorder()

		data := &testUser{Name: "Jane", Age: 25}

		err := NewSuccessResponse(w, http.StatusCreated, data)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		want := `{"data":{"name":"Jane","age":25}}`

		if strings.TrimSpace(w.Body.String()) != want {
			t.Errorf("got body %q, want %q", w.Body.String(), want)
		}
	})

	t.Run("Does Not Wrap Maps", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"foo": "bar"}

		err := NewSuccessResponse(w, http.StatusOK, data)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		want := `{"foo":"bar"}`

		if strings.TrimSpace(w.Body.String()) != want {
			t.Errorf("got body %q, want %q", w.Body.String(), want)
		}
	})

	t.Run("Does Not Wrap Slices", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []string{"a", "b"}

		err := NewSuccessResponse(w, http.StatusOK, data)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		want := `["a","b"]`

		if strings.TrimSpace(w.Body.String()) != want {
			t.Errorf("got body %q, want %q", w.Body.String(), want)
		}
	})
}

func TestNewErrorResponse(t *testing.T) {
	t.Run("Standard Error Response", func(t *testing.T) {
		w := httptest.NewRecorder()

		id := "test-id"
		msg := "Invalid input"
		details := []string{"field x is required"}

		err := NewErrorResponse(w, http.StatusBadRequest, msg, details, id)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
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
		err := NewErrorResponse(w, http.StatusNotFound, "Not Found", nil, "id")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		var response ErrorResponse

		err = json.Unmarshal(w.Body.Bytes(), &response)

		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Details == nil {
			t.Error("got nil, want empty slice for details")
		}

		if len(response.Error.Details) != 0 {
			t.Error("got non-zero length, want empty slice length 0")
		}
	})
}

func TestInternalNewResponse(t *testing.T) {
	t.Run("Write Failure (Network Error)", func(t *testing.T) {
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
		data := map[string]string{"foo": "bar"}

		err := newResponse(w, http.StatusOK, data)

		if err == nil {
			t.Fatal("got nil, want error from Write failure")
		}

		if !strings.Contains(err.Error(), "write response error") {
			t.Errorf("got %v, want 'write response error'", err)
		}
	})

	t.Run("Marshal Failure (Fallback Success)", func(t *testing.T) {
		w := httptest.NewRecorder()

		data := map[string]float64{"val": math.Inf(1)} // Infinity fails Marshal

		err := newResponse(w, http.StatusOK, data)

		if err == nil {
			t.Fatal("got nil, want error from marshal failure")
		}

		// It should contain the marshal error...
		if !strings.Contains(err.Error(), "marshal response error") {
			t.Fatalf("got %v, want 'marshal response error' in error", err)
		}

		// ...but NOT the fallback error (because fallback succeeded)
		if strings.Contains(err.Error(), "marshal fallback error") {
			t.Fatalf("got %v, want no 'marshal fallback error' in error", err)
		}

		// Verify it wrote the fallback 500 status
		if w.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want 500", w.Code)
		}
	})

	t.Run("Double Failure (Marshal + Write)", func(t *testing.T) {
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		// Infinity fails json.Marshal, triggering the fallback logic
		data := map[string]float64{"val": math.Inf(1)}

		err := newResponse(w, http.StatusOK, data)

		if err == nil {
			t.Fatal("got nil, want error from double failure")
		}

		// The error should contain both the marshal error AND the fallback write error
		if !strings.Contains(err.Error(), "marshal response error") || !strings.Contains(err.Error(), "marshal fallback error") {
			t.Errorf("got %v, want both 'marshal response error' and 'marshal fallback error'", err)
		}
	})
}
