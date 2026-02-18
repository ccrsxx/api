package jellyfin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ccrsxx/api/internal/clients/jellyfin"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/model"
)

type service struct {
	mu            sync.Mutex
	fetcher       func(context.Context) ([]jellyfin.SessionInfo, error)
	lastState     *model.CurrentlyPlaying
	lastStateTime time.Time
}

var Service = &service{
	fetcher: func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
		return jellyfin.DefaultClient().GetSessions(ctx)
	},
}

func (s *service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	sessions, err := s.fetcher(ctx)

	if err != nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin), fmt.Errorf("jellyfin get sessions error: %w", err)
	}

	var playingItem *model.CurrentlyPlaying

	for _, session := range sessions {
		isNotValidUsername := session.UserName == nil || *session.UserName != config.Env().JellyfinUsername

		if isNotValidUsername {
			continue
		}

		isNotPlayingSession := session.NowPlayingItem == nil || session.PlayState == nil

		if isNotPlayingSession {
			continue
		}

		isNotAudioSession := session.NowPlayingItem.Type != jellyfin.KindAudio

		if isNotAudioSession {
			continue
		}

		parsedValue := parseJellyfinSessions(&session)
		playingItem = &parsedValue

		break
	}

	if playingItem == nil {
		return s.getCachedStateOrEmpty(), nil
	}

	s.mu.Lock()

	s.lastState = playingItem
	s.lastStateTime = time.Now()

	s.mu.Unlock()

	return *playingItem, nil
}

func (s *service) getCachedStateOrEmpty() model.CurrentlyPlaying {
	s.mu.Lock()
	defer s.mu.Unlock()

	const gracePeriod = 5 * time.Second

	shouldUseCache := s.lastState != nil &&
		s.lastState.IsPlaying &&
		time.Since(s.lastStateTime) < gracePeriod

	if shouldUseCache {
		return s.getExtrapolatedState()
	}

	return model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin)
}

func (s *service) getExtrapolatedState() model.CurrentlyPlaying {
	if s.lastState == nil || s.lastState.Item == nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin)
	}

	elapsed := int(time.Since(s.lastStateTime).Milliseconds())
	extrapolatedProgress := s.lastState.Item.ProgressMs + elapsed

	progressMs := min(extrapolatedProgress, s.lastState.Item.DurationMs)

	itemCopy := *s.lastState.Item
	itemCopy.ProgressMs = progressMs

	return model.CurrentlyPlaying{
		Platform:  model.PlatformJellyfin,
		IsPlaying: s.lastState.IsPlaying,
		Item:      &itemCopy,
	}
}
