package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
)

func TestRecovery(t *testing.T) {
	// A handler that panics on purpose
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went terribly wrong")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	Recovery(panicHandler).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want 500", w.Code)
	}

	var res api.ErrorResponse

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	// In production (default config), the message should be generic
	if res.Error.Message != "An internal server error occurred" {
		t.Errorf("got message %q, want generic error message", res.Error.Message)
	}

	// 3. Ensure we didn't leak the panic message to the client (Security)
	if strings.Contains(w.Body.String(), "something went terribly wrong") {
		t.Error("panic message leaked to client body")
	}
}
