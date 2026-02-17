package jellyfin

import (
	"testing"

	"github.com/ccrsxx/api/src/internal/clients/jellyfin"
	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/model"
)

func Test_parseJellyfinSessions(t *testing.T) {
	// Setup config for image URL construction
	originalImgUrl := config.Env().JellyfinImageUrl

	defer func() {
		config.Env().JellyfinImageUrl = originalImgUrl
	}()

	config.Env().JellyfinImageUrl = "http://jellyfin.com"

	t.Run("Full Data", func(t *testing.T) {
		name := "Song"
		artist := "Artist"
		album := "Album"
		id := "item-1"

		ticks := int64(100000) // 10ms
		pos := int64(50000)    // 5ms

		session := &jellyfin.SessionInfo{
			NowPlayingItem: &jellyfin.BaseItem{
				Name:         &name,
				Artists:      []string{artist},
				Album:        &album,
				Id:           id,
				RunTimeTicks: &ticks,
			},
			PlayState: &jellyfin.PlayerStateInfo{
				IsPaused:      false,
				PositionTicks: &pos,
			},
		}

		got := parseJellyfinSessions(session)

		if got.Platform != model.PlatformJellyfin {
			t.Errorf("want platform jellyfin, got %s", got.Platform)
		}

		if !got.IsPlaying {
			t.Error("want playing true")
		}

		if got.Item.TrackName != "Song" {
			t.Errorf("want Song, got %s", got.Item.TrackName)
		}

		// 10000 ticks = 1ms. 50000 ticks = 5ms.
		if got.Item.ProgressMs != 5 {
			t.Errorf("want 5ms, got %d", got.Item.ProgressMs)
		}

		wantImg := "http://jellyfin.com/Items/item-1/Images/Primary"

		if *got.Item.AlbumImageURL != wantImg {
			t.Errorf("want %s, got %s", wantImg, *got.Item.AlbumImageURL)
		}
	})

	t.Run("Minimal Data (Fallbacks)", func(t *testing.T) {
		// Test logic where Name/Album are nil
		session := &jellyfin.SessionInfo{
			NowPlayingItem: &jellyfin.BaseItem{
				// All nil
			},
			PlayState: &jellyfin.PlayerStateInfo{
				IsPaused: true,
			},
		}

		got := parseJellyfinSessions(session)

		if got.IsPlaying {
			t.Fatal("want paused, got playing")
		}

		if got.Item.TrackName != "Unknown Track" {
			t.Fatalf("want fallback track name, got %s", got.Item.TrackName)
		}

		if got.Item.ArtistName != "Unknown Artist" {
			t.Fatalf("want fallback artist name, got %s", got.Item.ArtistName)
		}

		if got.Item.AlbumName != "Unknown Album" {
			t.Errorf("want fallback album name, got %s", got.Item.AlbumName)
		}
	})

	t.Run("Fallback to OriginalTitle", func(t *testing.T) {
		orig := "Original Album"

		session := &jellyfin.SessionInfo{
			NowPlayingItem: &jellyfin.BaseItem{
				OriginalTitle: &orig,
			},
			PlayState: &jellyfin.PlayerStateInfo{},
		}

		got := parseJellyfinSessions(session)

		if got.Item.AlbumName != "Original Album" {
			t.Errorf("want Original Album, got %s", got.Item.AlbumName)
		}
	})
}
