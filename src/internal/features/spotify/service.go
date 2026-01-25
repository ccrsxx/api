// src/features/spotify/service.go

package spotify

import (
	"context"
	"fmt"

	"github.com/ccrsxx/api-go/src/internal/clients/spotify"
	"github.com/ccrsxx/api-go/src/internal/model"
)

type service struct{}

var Service = &service{}

func (s *service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	data, err := spotify.Client().GetCurrentlyPlaying(ctx)

	if err != nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), fmt.Errorf("spotify get currently playing error: %w", err)
	}

	// Handle 204 No Content case
	if data == nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), nil
	}

	return parseSpotifyCurrentlyPlaying(data), nil
}
