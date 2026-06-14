package navidrome

import (
	"context"
	"fmt"
	"io"

	"github.com/ccrsxx/api/internal/clients/navidrome"
	"github.com/ccrsxx/api/internal/model"
)

type navidromeClient interface {
	GetNowPlaying(ctx context.Context) ([]navidrome.NowPlayingEntry, error)
	GetCoverArtStream(ctx context.Context, covertArtID string) (io.ReadCloser, error)
}

type Service struct {
	client            navidromeClient
	backendPublicURL  string
	navidromeUsername string
}

type ServiceConfig struct {
	Client            navidromeClient
	BackendPublicURL  string
	NavidromeUsername string
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		client:            cfg.Client,
		backendPublicURL:  cfg.BackendPublicURL,
		navidromeUsername: cfg.NavidromeUsername,
	}
}

func (s *Service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	nowPlaying, err := s.client.GetNowPlaying(ctx)

	if err != nil {
		return model.CurrentlyPlaying{}, fmt.Errorf("navidrome currently playing error: %w", err)
	}

	var playingItem *model.CurrentlyPlaying

	for _, entry := range nowPlaying {
		isNotValidUsername := entry.UserName != s.navidromeUsername

		if isNotValidUsername {
			continue
		}

		playingItem = new(parseNavidromeCurrentlyPlaying(entry, s.backendPublicURL))

		break
	}

	if playingItem == nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformNavidrome), nil
	}

	return *playingItem, nil

}
