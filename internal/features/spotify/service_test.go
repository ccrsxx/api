package spotify

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/clients/spotify"
)

func TestService_GetCurrentlyPlaying(t *testing.T) {
	originalFetcher := Service.fetcher

	defer func() {
		Service.fetcher = originalFetcher
	}()

	t.Run("Success Playing", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
			return &spotify.SpotifyCurrentlyPlaying{
				IsPlaying: true,
				Item: &spotify.SpotifyItem{
					Name: "Test Song",
					Album: &spotify.SpotifyAlbum{
						Name: "Test Album",
					},
				},
			}, nil
		}

		got, err := Service.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if !got.IsPlaying {
			t.Fatal("want playing true")
		}

		if got.Item.TrackName != "Test Song" {
			t.Errorf("got %s, want track Test Song", got.Item.TrackName)
		}
	})

	t.Run("Success Not Playing (Nil Data)", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
			return nil, nil // 204 No Content behavior
		}

		got, err := Service.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got.IsPlaying {
			t.Fatal("want playing false for nil data")
		}

		if got.Item != nil {
			t.Error("want nil item for idle state")
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
			return nil, errors.New("network fail")
		}

		got, err := Service.GetCurrentlyPlaying(context.Background())

		// The service WRAPS the error, so err should NOT be nil
		if err == nil {
			t.Error("want error, got nil")
		}

		// It should still return a default object even on error
		if got.IsPlaying {
			t.Error("want IsPlaying false on error")
		}
	})

	t.Run("Default Fetcher Execution for coverage", func(t *testing.T) {
		f := originalFetcher

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

		defer cancel()

		if _, err := f(ctx); err != nil {
			// We expect an error here since we likely don't have a Jellyfin server running during tests
			t.Logf("default fetcher returned want error (no server): %v", err)
			return
		}
	})
}
