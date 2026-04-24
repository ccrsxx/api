package home_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/features/home"
	"github.com/ccrsxx/api/internal/test"
)

func TestController_Ping(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.Host = "api.example.com"

		ctrl := home.NewController()

		ctrl.Ping(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}

		var res api.SuccessResponse[home.PingResponse]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		wantMsg := "Welcome to the API! The server is up and running."

		if res.Data.Message != wantMsg {
			t.Fatalf("got %q, want message %q", res.Data.Message, wantMsg)
		}

		wantDocURL := "http://api.example.com/docs"

		if res.Data.DocumentationURL != wantDocURL {
			t.Errorf("got %q, want docs url %q", res.Data.DocumentationURL, wantDocURL)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}

		ctrl := home.NewController()

		ctrl.Ping(w, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
