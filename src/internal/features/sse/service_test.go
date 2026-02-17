package sse

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/src/internal/model"
)

func TestService_IsConnectionAllowed(t *testing.T) {
	// Clean slate
	Service.clients = map[chan string]clientMetadata{}
	Service.ipAddressCounts = map[string]int{}

	t.Run("Allowed", func(t *testing.T) {
		if err := Service.IsConnectionAllowed("1.1.1.1"); err != nil {
			t.Errorf("want allowed, got %v", err)
		}
	})

	t.Run("IP Limit Reached", func(t *testing.T) {
		ip := "2.2.2.2"
		Service.ipAddressCounts[ip] = maxClientsPerIP

		err := Service.IsConnectionAllowed(ip)

		if err == nil {
			t.Error("want IP limit error, got nil")
		}
	})

	t.Run("Global Limit Reached", func(t *testing.T) {
		// Fake filling the map
		for range maxGlobalClients {
			c := make(chan string)
			Service.clients[c] = clientMetadata{}
		}

		err := Service.IsConnectionAllowed("3.3.3.3")
		if err == nil {
			t.Error("want global limit error, got nil")
		}

		// Cleanup
		Service.clients = map[chan string]clientMetadata{}
	})
}

func TestService_AddRemoveClient(t *testing.T) {
	// Mock Dependencies
	originalPoll := Service.pollInterval

	originalSpot := Service.spotifyFetcher
	originalJelly := Service.jellyfinFetcher

	defer func() {
		Service.pollInterval = originalPoll

		Service.spotifyFetcher = originalSpot
		Service.jellyfinFetcher = originalJelly
	}()

	Service.pollInterval = 10 * time.Millisecond

	Service.spotifyFetcher = func(ctx context.Context) (model.CurrentlyPlaying, error) {
		return model.CurrentlyPlaying{Platform: model.PlatformSpotify, IsPlaying: true}, nil
	}

	Service.jellyfinFetcher = func(ctx context.Context) (model.CurrentlyPlaying, error) {
		return model.CurrentlyPlaying{Platform: model.PlatformJellyfin, IsPlaying: false}, nil
	}

	t.Run("Add or Remove Client Lifecycle", func(t *testing.T) {
		// Use buffered channel matching production controller
		clientChan := make(chan string, 4)

		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		// Add Client
		Service.AddClient(ctx, clientChan, "127.0.0.1", "TestAgent")

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
		Service.RemoveClient(ctx, clientChan)

		if len(Service.clients) != 0 {
			t.Error("want clients map to be empty")
		}

		if Service.stopChan != nil {
			// Wait a bit for stopWorker to finish
			time.Sleep(20 * time.Millisecond)
			if Service.stopChan != nil {
				t.Error("want poller to stop")
			}
		}
	})

	t.Run("Client Channel Closed Before first message", func(t *testing.T) {
		clientChan := make(chan string, 4)

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		Service.AddClient(ctx, clientChan, "1.1.1.1", "TestAgent")

		Service.RemoveClient(ctx, clientChan)

		if len(Service.clients) != 0 {
			t.Error("want clients map to be empty after removing closed channel")
		}
	})
}

func TestService_getSSEData_Errors(t *testing.T) {
	// Mock errors
	Service.spotifyFetcher = func(ctx context.Context) (model.CurrentlyPlaying, error) {
		return model.CurrentlyPlaying{}, errors.New("spotify fail")
	}

	Service.jellyfinFetcher = func(ctx context.Context) (model.CurrentlyPlaying, error) {
		return model.CurrentlyPlaying{}, errors.New("jellyfin fail")
	}

	// Should not panic, should return default empty structs
	data := Service.getSSEData(context.Background())

	if !strings.Contains(data.spotify, "spotify") {
		t.Error("want spotify event structure even on error")
	}
}

func TestService_WorkerLocks(t *testing.T) {
	Service.stopChan = nil

	Service.startWorkerLocked()
	if Service.stopChan == nil {
		t.Error("want worker started")
	}

	Service.startWorkerLocked()

	Service.stopWorkerLocked()

	if Service.stopChan != nil {
		t.Error("want worker stopped")
	}

	Service.stopWorkerLocked()
}
