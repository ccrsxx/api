package spotify

import (
	"testing"

	"github.com/ccrsxx/api/src/internal/clients/spotify"
	"github.com/ccrsxx/api/src/internal/model"
)

func Test_parseSpotifyCurrentlyPlaying(t *testing.T) {
	t.Run("Full Data", func(t *testing.T) {
		raw := &spotify.SpotifyCurrentlyPlaying{
			IsPlaying:  true,
			ProgressMs: 1000,
			Item: &spotify.SpotifyItem{
				Name:       "Track",
				DurationMs: 3000,
				IsLocal:    false,
				Artists: []spotify.SpotifyArtist{
					{Name: "A1"}, {Name: "A2"},
				},
				Album: &spotify.SpotifyAlbum{
					Name: "Album",
					Images: []spotify.SpotifyImage{
						{URL: "http://img.com/1"},
					},
				},
				ExternalURLs: spotify.SpotifyExternalURLs{
					Spotify: "http://open.&spotify.com/track/1",
				},
			},
		}

		got := parseSpotifyCurrentlyPlaying(raw)

		if got.Platform != model.PlatformSpotify {
			t.Errorf("want platform spotify, got %s", got.Platform)
		}

		if got.Item.ArtistName != "A1, A2" {
			t.Errorf("want 'A1, A2', got %q", got.Item.ArtistName)
		}

		if *got.Item.TrackURL != "http://open.&spotify.com/track/1" {
			t.Error("wrong track url")
		}

		if *got.Item.AlbumImageURL != "http://img.com/1" {
			t.Error("wrong image url")
		}
	})

	t.Run("Local File (No URLs)", func(t *testing.T) {
		raw := &spotify.SpotifyCurrentlyPlaying{
			IsPlaying: true,
			Item: &spotify.SpotifyItem{
				Name:    "Local Song",
				IsLocal: true,
				Album: &spotify.SpotifyAlbum{
					Images: []spotify.SpotifyImage{{URL: "should-ignore"}},
				},
				ExternalURLs: spotify.SpotifyExternalURLs{
					Spotify: "should-ignore",
				},
			},
		}

		got := parseSpotifyCurrentlyPlaying(raw)

		if got.Item.TrackURL != nil {
			t.Error("want nil track url for local file")
		}

		if got.Item.AlbumImageURL != nil {
			t.Error("want nil image url for local file")
		}
	})

	t.Run("No Images", func(t *testing.T) {
		raw := &spotify.SpotifyCurrentlyPlaying{
			Item: &spotify.SpotifyItem{
				Name: "Track",
				Album: &spotify.SpotifyAlbum{
					Images: []spotify.SpotifyImage{},
				},
			},
		}

		got := parseSpotifyCurrentlyPlaying(raw)

		if got.Item.AlbumImageURL != nil {
			t.Error("want nil image url")
		}
	})

	t.Run("Nil Item", func(t *testing.T) {
		raw := &spotify.SpotifyCurrentlyPlaying{
			IsPlaying: true,
			Item:      nil,
		}

		got := parseSpotifyCurrentlyPlaying(raw)

		if got.Item != nil {
			t.Fatalf("want nil item when raw.Item is nil, got %+v", got.Item)
		}

		if !got.IsPlaying {
			t.Error("want IsPlaying to be preserved")
		}
	})
}
