package sse

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/model"
)

type mockDataFetcher struct {
	result model.CurrentlyPlaying
	err    error
}

func (m *mockDataFetcher) GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error) {
	return m.result, m.err
}

func TestService_IsConnectionAllowed(t *testing.T) {
	t.Run("Allowed", func(t *testing.T) {
		svc := NewService(ServiceConfig{})

		if err := svc.IsConnectionAllowed("1.1.1.1"); err != nil {
			t.Errorf("got %v, want allowed", err)
		}
	})

	t.Run("IP Limit Reached", func(t *testing.T) {
		svc := NewService(ServiceConfig{})
		ip := "2.2.2.2"

		// Access private fields safely since we are in the sse package
		svc.ipAddressCounts[ip] = maxClientsPerIP

		err := svc.IsConnectionAllowed(ip)

		if err == nil {
			t.Error("want IP limit error, got nil")
		}
	})

	t.Run("Global Limit Reached", func(t *testing.T) {
		svc := NewService(ServiceConfig{})

		// Fake filling the map
		for range maxGlobalClients {
			c := make(chan string)
			svc.clients[c] = clientMetadata{}
		}

		err := svc.IsConnectionAllowed("3.3.3.3")

		if err == nil {
			t.Error("want global limit error, got nil")
		}
	})
}

func TestService_AddRemoveClient(t *testing.T) {
	setupService := func() *Service {
		dummySpotify := &mockDataFetcher{
			result: model.NewDefaultCurrentlyPlaying(model.PlatformSpotify),
		}

		dummyJellyfin := &mockDataFetcher{
			result: model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin),
		}

		return NewService(ServiceConfig{
			PollInterval:    10 * time.Millisecond,
			SpotifyService:  dummySpotify,
			JellyfinService: dummyJellyfin,
		})
	}

	t.Run("Add or Remove Client Lifecycle", func(t *testing.T) {
		svc := setupService()

		// Use buffered channel matching production controller
		clientChan := make(chan string, 4)

		ctx := t.Context()

		// Add Client
		svc.AddClient(ctx, clientChan, "127.0.0.1", "TestAgent")

		// Verify Initial Data (Welcome + Spotify + Jellyfin)
		timeout := time.After(1 * time.Second)

		msgCount := 0

		// We expect 3 initial messages
		for i := range 3 {
			select {
			case <-timeout:
				t.Fatal("timeout waiting for initial messages")
			case msg := <-clientChan:
				msgCount++

				if i == 0 && !strings.Contains(msg, "Connection established") {
					t.Error("want welcome message first")
				}
			}
		}

		// Verify Polling Broadcast (Wait for next tick)
		select {
		case msg := <-clientChan:
			// Should receive updates from poll loop
			if !strings.Contains(msg, "event: spotify") && !strings.Contains(msg, "event: jellyfin") {
				t.Errorf("unwanted broadcast message: %s", msg)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for broadcast")
		}

		// Remove Client
		svc.RemoveClient(ctx, clientChan)

		if len(svc.clients) != 0 {
			t.Error("want clients map to be empty")
		}

		if svc.stopChan != nil {
			// Wait a bit for stopWorker to finish
			time.Sleep(10 * time.Millisecond)

			if svc.stopChan != nil {
				t.Error("want poller to stop")
			}
		}
	})

	t.Run("Client Channel Closed Before first message", func(t *testing.T) {
		svc := setupService()

		clientChan := make(chan string, 4)

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		svc.AddClient(ctx, clientChan, "1.1.1.1", "TestAgent")
		svc.RemoveClient(ctx, clientChan)

		if len(svc.clients) != 0 {
			t.Error("want clients map to be empty after removing closed channel")
		}
	})
}

func TestService_getSSEData_Errors(t *testing.T) {
	failSpotify := &mockDataFetcher{
		err: errors.New("spotify fail"),
	}

	failJellyfin := &mockDataFetcher{
		err: errors.New("jellyfin fail"),
	}

	svc := NewService(ServiceConfig{
		SpotifyService:  failSpotify,
		JellyfinService: failJellyfin,
	})

	// Should not panic, should return default empty structs
	data := svc.getSSEData(context.Background())

	if !strings.Contains(data.spotify, "spotify") {
		t.Error("want spotify event structure even on error")
	}
}

func TestService_WorkerLocks(t *testing.T) {
	svc := NewService(ServiceConfig{})

	if svc.stopChan != nil {
		t.Error("want initial stopChan to be nil")
	}

	svc.startWorkerLocked()

	if svc.stopChan == nil {
		t.Error("want worker started")
	}

	// Starting again shouldn't overwrite or panic
	svc.startWorkerLocked()

	svc.stopWorkerLocked()

	if svc.stopChan != nil {
		t.Error("want worker stopped")
	}

	// Stopping again shouldn't panic
	svc.stopWorkerLocked()
}
