package navidrome

import (
	"testing"

	"github.com/ccrsxx/api/internal/clients/navidrome"
	"github.com/ccrsxx/api/internal/model"
)

func Test_parseNavidromeCurrentlyPlaying(t *testing.T) {
	publicURL := "http://api.example.com"

	t.Run("Playing", func(t *testing.T) {
		entry := navidrome.NowPlayingEntry{
			Child: navidrome.Child{
				Title:    "Song",
				Album:    "Album",
				Artist:   "Artist",
				CoverArt: "ca-123",
				Duration: 300,
			},
			UserName:   "user1",
			State:      "playing",
			PositionMs: 5000,
		}

		got := parseNavidromeCurrentlyPlaying(entry, publicURL)

		if got.Platform != model.PlatformNavidrome {
			t.Errorf("got %s, want platform navidrome", got.Platform)
		}

		if !got.IsPlaying {
			t.Error("want playing true")
		}

		if got.Item.TrackName != "Song" {
			t.Errorf("got %s, want Song", got.Item.TrackName)
		}

		if got.Item.AlbumName != "Album" {
			t.Errorf("got %s, want Album", got.Item.AlbumName)
		}

		if got.Item.ArtistName != "Artist" {
			t.Errorf("got %s, want Artist", got.Item.ArtistName)
		}

		if got.Item.ProgressMs != 5000 {
			t.Errorf("got %d, want 5000", got.Item.ProgressMs)
		}

		// Duration is in seconds, converted to ms (* 1000)
		if got.Item.DurationMs != 300000 {
			t.Errorf("got %d, want 300000", got.Item.DurationMs)
		}

		wantImg := "http://api.example.com/navidrome/cover-art/ca-123"

		if *got.Item.AlbumImageURL != wantImg {
			t.Errorf("got %s, want %s", *got.Item.AlbumImageURL, wantImg)
		}

		if got.Item.TrackURL != nil {
			t.Error("want nil TrackURL for navidrome")
		}
	})

	t.Run("Paused", func(t *testing.T) {
		entry := navidrome.NowPlayingEntry{
			State: "paused",
		}
		got := parseNavidromeCurrentlyPlaying(entry, publicURL)
		if got.IsPlaying {
			t.Error("want playing false for paused state")
		}
	})
}
