// src/features/spotify/service.go

package spotify

import (
	"context"
	"fmt"

	"github.com/ccrsxx/api-go/src/internal/clients/spotify"
)

type service struct{}

var Service = &service{}

func (s *service) GetCurrentlyPlaying(ctx context.Context) (*CurrentlyPlaying, error) {
	data, err := spotify.Client().GetNowCurrentlyPlaying(ctx)

	if err != nil {
		return nil, fmt.Errorf("spotify get currently playing error: %w", err)
	}

	if data == nil {
		return getDefaultCurrentlyPlaying(), nil
	}

	return mapSpotifyCurrentlyPlaying(data), nil
}
