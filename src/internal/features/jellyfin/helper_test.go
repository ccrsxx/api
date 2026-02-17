package jellyfin

import (
	"testing"

	"github.com/ccrsxx/api/src/internal/clients/jellyfin"
	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/model"
)

func Test_parseJellyfinSessions(t *testing.T) {
	originalImgUrl := config.Env().JellyfinImageUrl

	defer func() {
		config.Env().JellyfinImageUrl = originalImgUrl
	}()

	config.Env().JellyfinImageUrl = "http://jellyfin.com"

	t.Run("Full Data", func(t *testing.T) {
		id := "item-1"
		name := "Song"
		album := "Album"
		artist := "Artist"

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
			t.Errorf("got %s, want platform jellyfin", got.Platform)
		}

		if !got.IsPlaying {
			t.Error("want playing true")
		}

		if got.Item.TrackName != "Song" {
			t.Errorf("got %s, want Song", got.Item.TrackName)
		}

		// 10000 ticks = 1ms. 50000 ticks = 5ms.
		if got.Item.ProgressMs != 5 {
			t.Errorf("got %d, want 5ms", got.Item.ProgressMs)
		}

		wantImg := "http://jellyfin.com/Items/item-1/Images/Primary"

		if *got.Item.AlbumImageURL != wantImg {
			t.Errorf("got %s, want %s", *got.Item.AlbumImageURL, wantImg)
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
			t.Fatal("got playing, want paused")
		}

		if got.Item.TrackName != "Unknown Track" {
			t.Fatalf("got %s, want fallback track name", got.Item.TrackName)
		}

		if got.Item.ArtistName != "Unknown Artist" {
			t.Fatalf("got %s, want fallback artist name", got.Item.ArtistName)
		}

		if got.Item.AlbumName != "Unknown Album" {
			t.Errorf("got %s, want fallback album name", got.Item.AlbumName)
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
			t.Errorf("got %s, want Original Album", got.Item.AlbumName)
		}
	})
}
