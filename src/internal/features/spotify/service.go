package spotify

import (
	"context"
	"fmt"

	"github.com/ccrsxx/api/src/internal/clients/spotify"
	"github.com/ccrsxx/api/src/internal/model"
)

type service struct {
	fetcher func(context.Context) (*spotify.SpotifyCurrentlyPlaying, error)
}

var Service = &service{
	fetcher: func(ctx context.Context) (*spotify.SpotifyCurrentlyPlaying, error) {
		return spotify.DefaultClient().GetCurrentlyPlaying(ctx)
	},
}

func (s *service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	data, err := s.fetcher(ctx)

	if err != nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), fmt.Errorf("spotify get currently playing error: %w", err)
	}

	// Handle 204 No Content case
	if data == nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), nil
	}

	return parseSpotifyCurrentlyPlaying(data), nil
}
