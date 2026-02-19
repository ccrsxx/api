package jellyfin

import (
	"fmt"
	"strings"

	"github.com/ccrsxx/api/internal/clients/jellyfin"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/model"
)

func parseJellyfinSessions(session jellyfin.SessionInfo) model.CurrentlyPlaying {
	item := session.NowPlayingItem
	playState := session.PlayState

	trackName := "Unknown Track"

	if item.Name != nil {
		trackName = *item.Name
	}

	artistName := "Unknown Artist"

	if len(item.Artists) > 0 {
		artistName = strings.Join(item.Artists, ", ")
	}

	albumName := "Unknown Album"

	if item.Album != nil {
		albumName = *item.Album
	} else if item.OriginalTitle != nil {
		albumName = *item.OriginalTitle
	}

	albumImageUrl := ""

	if item.Id != "" {
		albumImageUrl = fmt.Sprintf("%s/Items/%s/Images/Primary", config.Env().JellyfinImageUrl, item.Id)
	}

	// Jellyfin uses "Ticks". 1 ms = 10,000 Ticks. Convert to ms.

	durationMs := 0

	if item.RunTimeTicks != nil {
		durationMs = int(*item.RunTimeTicks / 10000)
	}

	progressMs := 0

	if playState.PositionTicks != nil {
		progressMs = int(*playState.PositionTicks / 10000)
	}

	return model.CurrentlyPlaying{
		Platform:  model.PlatformJellyfin,
		IsPlaying: !playState.IsPaused,
		Item: &model.Track{
			TrackURL:      nil, // Jellyfin has no public track URL
			TrackName:     trackName,
			AlbumName:     albumName,
			ArtistName:    artistName,
			DurationMs:    durationMs,
			ProgressMs:    progressMs,
			AlbumImageURL: &albumImageUrl,
		},
	}
}
