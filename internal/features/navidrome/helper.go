package navidrome

import (
	"fmt"

	"github.com/ccrsxx/api/internal/clients/navidrome"
	"github.com/ccrsxx/api/internal/model"
)

func parseNavidromeCurrentlyPlaying(raw navidrome.NowPlayingEntry, publicURL string) model.CurrentlyPlaying {
	isPlaying := raw.State == "playing"
	albumImageURL := fmt.Sprintf("%s/navidrome/cover-art/%s", publicURL, raw.CoverArt)

	return model.CurrentlyPlaying{
		Platform:  model.PlatformNavidrome,
		IsPlaying: isPlaying,
		Item: &model.Track{
			TrackURL:      nil, // Navidrome has no public track URL
			TrackName:     raw.Title,
			AlbumName:     raw.Album,
			ArtistName:    raw.Artist,
			ProgressMs:    int(raw.PositionMs),
			DurationMs:    int(raw.Duration) * 1000,
			AlbumImageURL: new(albumImageURL),
		},
	}
}
