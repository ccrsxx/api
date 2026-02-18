package jellyfin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/jellyfin"
	"github.com/ccrsxx/api/internal/test"
)

func TestController_getCurrentlyPlaying(t *testing.T) {
	originalFetcher := Service.fetcher

	defer func() {
		Service.fetcher = originalFetcher
	}()

	t.Run("Success", func(t *testing.T) {
		// We return nil sessions, which results in "Not Playing" (200 OK)
		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return nil, nil
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		Controller.getCurrentlyPlaying(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res api.SuccessResponse[struct {
			IsPlaying bool `json:"isPlaying"`
		}]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res.Data.IsPlaying {
			t.Error("want isPlaying false")
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return nil, errors.New("fail")
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		Controller.getCurrentlyPlaying(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want 500, got %d", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return nil, nil
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		Controller.getCurrentlyPlaying(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
