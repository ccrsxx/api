package jellyfin

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ccrsxx/api/src/internal/clients/jellyfin"
	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/model"
)

func TestService_GetCurrentlyPlaying(t *testing.T) {
	originalFetcher := Service.fetcher

	originalUser := config.Env().JellyfinUsername

	defer func() {
		Service.fetcher = originalFetcher

		config.Env().JellyfinUsername = originalUser
	}()

	config.Env().JellyfinUsername = "testuser"

	validUser := "testuser"
	otherUser := "other"

	t.Run("Success Playing", func(t *testing.T) {
		resetJellyfinCache()

		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return []jellyfin.SessionInfo{
				{
					UserName: &validUser,
					NowPlayingItem: &jellyfin.BaseItem{
						Name: new("Song"),
						Type: jellyfin.KindAudio,
					},
					PlayState: &jellyfin.PlayerStateInfo{IsPaused: false},
				},
			}, nil
		}

		got, err := Service.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if !got.IsPlaying {
			t.Error("want playing true")
		}

		if got.Item.TrackName != "Song" {
			t.Errorf("got %s, want Song", got.Item.TrackName)
		}
	})

	t.Run("Filtering Logic", func(t *testing.T) {
		resetJellyfinCache()

		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return []jellyfin.SessionInfo{
				{UserName: &otherUser},                      // Skip: Wrong User
				{UserName: &validUser, NowPlayingItem: nil}, // Skip: Not Playing
				{
					UserName:       &validUser,
					NowPlayingItem: &jellyfin.BaseItem{Type: jellyfin.KindMovie}, // Skip: Not Audio
					PlayState:      &jellyfin.PlayerStateInfo{},
				},
			}, nil
		}

		got, err := Service.GetCurrentlyPlaying(context.Background())

		fmt.Printf("Got state: %+v\n", got)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got.IsPlaying {
			t.Error("want not playing (all sessions filtered)")
		}
	})

	t.Run("Fetcher Error", func(t *testing.T) {
		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return nil, errors.New("network fail")
		}

		_, err := Service.GetCurrentlyPlaying(context.Background())
		if err == nil {
			t.Error("want error")
		}
	})

	t.Run("Caching and Extrapolation", func(t *testing.T) {
		// Prime the cache with a playing state
		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return []jellyfin.SessionInfo{
				{
					UserName: &validUser,
					NowPlayingItem: &jellyfin.BaseItem{
						Name:         new("Cached Song"),
						Type:         jellyfin.KindAudio,
						RunTimeTicks: new(int64(60000000)), // 6000ms duration (using large numbers for ticks)
					},
					PlayState: &jellyfin.PlayerStateInfo{
						IsPaused:      false,
						PositionTicks: new(int64(10000000)), // 1000ms progress
					},
				},
			}, nil
		}

		_, err := Service.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		// Simulate "No Content" from API (user momentarily between songs or network blip)
		// But within grace period
		Service.fetcher = func(ctx context.Context) ([]jellyfin.SessionInfo, error) {
			return []jellyfin.SessionInfo{}, nil
		}

		// Force time to advance slightly for extrapolation
		// Note: We can't easily mock time.Now() globally, but we can verify
		// that the logic *attempts* to return the cached item.
		// For robustness, we manually tweak the lastStateTime in the struct if needed,
		// but since we just ran it, it is definitely within the 5s grace period.

		got, err := Service.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if !got.IsPlaying {
			t.Error("want cached playing state")
		}

		if got.Item.TrackName != "Cached Song" {
			t.Error("want cached song data")
		}

		// Force Cache Expiry
		// Manually reach into the service and set time back > 5 seconds
		Service.mu.Lock()

		Service.lastStateTime = time.Now().Add(-10 * time.Second)

		Service.mu.Unlock()

		gotExpired, err := Service.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if gotExpired.IsPlaying {
			t.Error("want empty state after cache expiry")
		}
	})

	t.Run("Extrapolation Logic Edge Cases", func(t *testing.T) {
		// Manually setup a state where Progress is near Duration
		Service.lastState = &model.CurrentlyPlaying{
			IsPlaying: true,
			Item: &model.Track{
				ProgressMs: 5000,
				DurationMs: 5100,
			},
		}

		// Set lastStateTime to 200ms ago, which would extrapolate progress to 5200ms (beyond duration)
		Service.lastStateTime = time.Now().Add(-200 * time.Millisecond)

		// Call internal method directly to test math
		got := Service.getExtrapolatedState()

		// Should have advanced by ~200ms, but capped at 5100
		if got.Item.ProgressMs > 5100 {
			t.Errorf("got %d, want capped progress", got.Item.ProgressMs)
		}

		// Test nil checks in getExtrapolatedState
		Service.lastState = nil

		gotEmpty := Service.getExtrapolatedState()

		if gotEmpty.IsPlaying {
			t.Error("want empty state for nil lastState")
		}
	})

	t.Run("Default Fetcher Execution for coverage", func(t *testing.T) {
		f := originalFetcher

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

		defer cancel()

		if _, err := f(ctx); err != nil {
			// We expect an error here since we likely don't have a Jellyfin server running during tests
			t.Logf("default fetcher returned want error (no server): %v", err)
			return
		}
	})
}

func resetJellyfinCache() {
	Service.lastState = nil
}
