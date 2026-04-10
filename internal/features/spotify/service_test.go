package spotify

import (
	"context"
	"errors"
	"testing"

	"github.com/ccrsxx/api/internal/clients/spotify"
)

func TestService_GetCurrentlyPlaying(t *testing.T) {
	t.Run("Success Playing", func(t *testing.T) {
		mockFetcher := func(ctx context.Context) (spotify.SpotifyCurrentlyPlaying, error) {
			return spotify.SpotifyCurrentlyPlaying{
				IsPlaying: true,
				Item: &spotify.SpotifyItem{
					Name: "Test Song",
					Album: &spotify.SpotifyAlbum{
						Name: "Test Album",
					},
				},
			}, nil
		}

		svc := NewService(ServiceConfig{
			Fetcher: mockFetcher,
		})

		got, err := svc.GetCurrentlyPlaying(context.Background())

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

	t.Run("Success No Content", func(t *testing.T) {
		mockFetcher := func(ctx context.Context) (spotify.SpotifyCurrentlyPlaying, error) {
			return spotify.SpotifyCurrentlyPlaying{}, spotify.ErrNoContent
		}

		svc := NewService(ServiceConfig{
			Fetcher: mockFetcher,
		})

		got, err := svc.GetCurrentlyPlaying(context.Background())

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
		mockFetcher := func(ctx context.Context) (spotify.SpotifyCurrentlyPlaying, error) {
			return spotify.SpotifyCurrentlyPlaying{}, errors.New("network fail")
		}

		svc := NewService(ServiceConfig{
			Fetcher: mockFetcher,
		})

		got, err := svc.GetCurrentlyPlaying(context.Background())

		// The service WRAPS the error, so err should NOT be nil
		if err == nil {
			t.Error("want error, got nil")
		}

		// It should still return a default object even on error
		if got.IsPlaying {
			t.Error("want IsPlaying false on error")
		}
	})
}
