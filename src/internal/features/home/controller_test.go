package home

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/api"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestController_ping(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.Host = "api.example.com"

		Controller.ping(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want status 200, got %d", w.Code)
		}

		var res api.SuccessResponse[struct {
			Message          string `json:"message"`
			DocumentationURL string `json:"documentationUrl"`
		}]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		wantMsg := "Welcome to the API! The server is up and running."

		if res.Data.Message != wantMsg {
			t.Fatalf("want message %q, got %q", wantMsg, res.Data.Message)
		}

		wantDocUrl := "http://api.example.com/docs"

		if res.Data.DocumentationURL != wantDocUrl {
			t.Errorf("want docs url %q, got %q", wantDocUrl, res.Data.DocumentationURL)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		Controller.ping(w, r)
	})
}
