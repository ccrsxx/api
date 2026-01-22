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

type service struct {
	*spotifyService
}

var Service = &service{
	spotifyService: &spotifyService{},
}

type spotifyService struct {
	mu            sync.Mutex
	lastState     *model.CurrentlyPlaying
	lastStateTime time.Time
}

func (ss *spotifyService) GetCurrentlyPlaying(ctx context.Context, useCache bool) (*model.CurrentlyPlaying, error) {
	sessions, err := jellyfin.Client().GetSessions(ctx)

	if err != nil {
		return nil, fmt.Errorf("jellyfin get sessions error: %w", err)
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

		playingItem = parseJellyfinSessions(&session)

		break
	}

	if playingItem == nil {
		return ss.getCachedStateOrEmpty(), nil
	}

	ss.mu.Lock()

	ss.lastState = playingItem
	ss.lastStateTime = time.Now()

	ss.mu.Unlock()

	return playingItem, nil
}

func (ss *spotifyService) getCachedStateOrEmpty() *model.CurrentlyPlaying {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	const gracePeriod = 5 * time.Second

	shouldUseCache := ss.lastState != nil &&
		ss.lastState.IsPlaying &&
		time.Since(ss.lastStateTime) < gracePeriod

	if shouldUseCache {
		return ss.getExtrapolatedState()
	}

	return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify)
}

func (ss *spotifyService) getExtrapolatedState() *model.CurrentlyPlaying {
	if ss.lastState == nil || ss.lastState.Item == nil {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify)
	}

	elapsed := int(time.Since(ss.lastStateTime).Milliseconds())
	extrapolatedProgress := ss.lastState.Item.ProgressMs + elapsed

	progressMs := min(extrapolatedProgress, ss.lastState.Item.DurationMs)

	itemCopy := *ss.lastState.Item
	itemCopy.ProgressMs = progressMs

	return &model.CurrentlyPlaying{
		Platform:  model.PlatformSpotify,
		IsPlaying: ss.lastState.IsPlaying,
		Item:      &itemCopy,
	}
}
