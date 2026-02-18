package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/test"
)

func TestController_getCurrentlyPlaying(t *testing.T) {
	originalFetcher := Service.fetcher

	defer func() {
		Service.fetcher = originalFetcher
	}()

	t.Run("Success", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
			return &spotify.SpotifyCurrentlyPlaying{
				IsPlaying: true,
				Item:      &spotify.SpotifyItem{Name: "Song"},
			}, nil
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		Controller.getCurrentlyPlaying(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[struct {
			IsPlaying bool `json:"isPlaying"`
		}]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if !res.Data.IsPlaying {
			t.Error("want isPlaying true")
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
			return nil, errors.New("fail")
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		Controller.getCurrentlyPlaying(w, r)

		// Service error returns the error to the controller, which calls HandleHttpError
		// Since it's a generic error, it usually results in 500
		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
			return nil, nil
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		Controller.getCurrentlyPlaying(errWriter, r)
	})
}
