package spotify

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/model"
)

type service struct {
	fetcher func(context.Context) (spotify.SpotifyCurrentlyPlaying, error)
}

type Config struct {
	Fetcher func(context.Context) (spotify.SpotifyCurrentlyPlaying, error)
}

func NewService(cfg Config) *service {
	return &service{
		fetcher: cfg.Fetcher,
	}
}

func (s *service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	data, err := s.fetcher(ctx)

	// Handle 204 No Content case
	if errors.Is(err, spotify.ErrNoContent) {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), nil
	}

	if err != nil {
		return model.CurrentlyPlaying{}, fmt.Errorf("spotify get currently playing error: %w", err)
	}

	return parseSpotifyCurrentlyPlaying(data), nil
}
