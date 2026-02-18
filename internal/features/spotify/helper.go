package spotify

import (
	"strings"

	"github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/model"
)

func parseSpotifyCurrentlyPlaying(raw *spotify.SpotifyCurrentlyPlaying) model.CurrentlyPlaying {
	item := raw.Item

	var artistNames []string

	for _, artist := range item.Artists {
		artistNames = append(artistNames, artist.Name)
	}

	artistStr := strings.Join(artistNames, ", ")

	var trackUrl *string
	var imageUrl *string

	albumName := "Unknown Album"

	if !item.IsLocal {
		// Only set URL if it exists
		if url := item.ExternalURLs.Spotify; url != "" {
			trackUrl = &url
		}

		// Only set album info if it exists
		if item.Album != nil {
			albumName = item.Album.Name

			if len(item.Album.Images) > 0 {
				url := item.Album.Images[0].URL
				imageUrl = &url
			}
		}
	}

	return model.CurrentlyPlaying{
		Platform:  model.PlatformSpotify,
		IsPlaying: raw.IsPlaying,
		Item: &model.Track{
			TrackURL:      trackUrl,
			TrackName:     item.Name,
			ArtistName:    artistStr,
			ProgressMs:    raw.ProgressMs,
			DurationMs:    item.DurationMs,
			AlbumName:     albumName,
			AlbumImageURL: imageUrl,
		},
	}
}
