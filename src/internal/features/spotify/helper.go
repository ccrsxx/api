package spotify

import (
	"strings"

	"github.com/ccrsxx/api-go/src/internal/clients/spotify"
	"github.com/ccrsxx/api-go/src/internal/model"
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

	if !item.IsLocal {
		// Only set URL if it exists
		if url := item.ExternalURLs.Spotify; url != "" {
			trackUrl = &url
		}

		// Get the first image (usually the largest)
		if len(item.Album.Images) > 0 {
			url := item.Album.Images[0].URL
			imageUrl = &url
		}
	}

	return model.CurrentlyPlaying{
		Platform:  model.PlatformSpotify,
		IsPlaying: raw.IsPlaying,
		Item: &model.Track{
			TrackName:     item.Name,
			AlbumName:     item.Album.Name,
			ArtistName:    artistStr,
			ProgressMs:    raw.ProgressMs,
			DurationMs:    item.DurationMs,
			TrackURL:      trackUrl, // string | null
			AlbumImageURL: imageUrl, // string | null
		},
	}
}
