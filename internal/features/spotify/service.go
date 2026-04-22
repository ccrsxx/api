package spotify

import (
	"context"
	"errors"
	"fmt"

	"github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/model"
)

type spotifyClient interface {
	GetCurrentlyPlaying(context.Context) (spotify.SpotifyCurrentlyPlaying, error)
}

type Service struct {
	client spotifyClient
}

type ServiceConfig struct {
	Client spotifyClient
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		client: cfg.Client,
	}
}

func (s *Service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	data, err := s.client.GetCurrentlyPlaying(ctx)

	// Handle 204 No Content case
	if errors.Is(err, spotify.ErrNoContent) {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), nil
	}

	if err != nil {
		return model.CurrentlyPlaying{}, fmt.Errorf("spotify get currently playing error: %w", err)
	}

	return parseSpotifyCurrentlyPlaying(data), nil
}
