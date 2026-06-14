package navidrome_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ccrsxx/api/internal/clients/navidrome"
	"github.com/ccrsxx/api/internal/model"

	navidromeFeature "github.com/ccrsxx/api/internal/features/navidrome"
)

func TestService_GetCurrentlyPlaying(t *testing.T) {
	validUser := "testuser"

	t.Run("Success Playing", func(t *testing.T) {
		mock := &mockNavidromeClient{
			nowPlayingResult: []navidrome.NowPlayingEntry{
				{
					Child: navidrome.Child{
						Title: "Song",
					},
					UserName: validUser,
					State:    "playing",
				},
			},
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client:           mock,
			NavidromeUsername: validUser,
		})

		got, err := svc.GetCurrentlyPlaying(context.Background())

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

	t.Run("Filtering Wrong Username", func(t *testing.T) {
		mock := &mockNavidromeClient{
			nowPlayingResult: []navidrome.NowPlayingEntry{
				{
					UserName: "other",
					State:    "playing",
				},
			},
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client:           mock,
			NavidromeUsername: validUser,
		})

		got, err := svc.GetCurrentlyPlaying(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if got.IsPlaying {
			t.Error("want not playing (wrong username filtered)")
		}

		if got.Platform != model.PlatformNavidrome {
			t.Errorf("got %s, want navidrome platform", got.Platform)
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mock := &mockNavidromeClient{
			nowPlayingErr: errors.New("network fail"),
		}

		svc := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
			Client: mock,
		})

		_, err := svc.GetCurrentlyPlaying(context.Background())

		if err == nil {
			t.Error("want error")
		}
	})
}
