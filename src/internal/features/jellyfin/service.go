package jellyfin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ccrsxx/api-go/src/internal/clients/jellyfin"
	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ccrsxx/api-go/src/internal/model"
)

var Service = &service{}

type service struct {
	mu            sync.Mutex
	lastState     *model.CurrentlyPlaying
	lastStateTime time.Time
}

func (s *service) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	sessions, err := jellyfin.Client().GetSessions(ctx)

	if err != nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin), fmt.Errorf("jellyfin get sessions error: %w", err)
	}

	var playingItem *model.CurrentlyPlaying

	for _, session := range *sessions {
		isValidUsername := session.UserName != nil || *session.UserName != config.Env().JellyfinUsername

		if !isValidUsername {
			continue
		}

		isNonPlayingSession := session.NowPlayingItem == nil || session.PlayState == nil

		if isNonPlayingSession {
			continue
		}

		isNonAudioSession := session.NowPlayingItem.Type != jellyfin.KindAudio

		if isNonAudioSession {
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
