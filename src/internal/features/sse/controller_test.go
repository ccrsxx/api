package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccrsxx/api/src/internal/model"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestController_getCurrentPlayingSSE(t *testing.T) {
	// Setup Mocks
	originalPollInterval := Service.pollInterval
	originalSpotifyFetcher := Service.spotifyFetcher
	originalJellyfinFetcher := Service.jellyfinFetcher

	defer func() {
		Service.pollInterval = originalPollInterval
		Service.spotifyFetcher = originalSpotifyFetcher
		Service.jellyfinFetcher = originalJellyfinFetcher
	}()

	Service.spotifyFetcher = func(ctx context.Context) (model.CurrentlyPlaying, error) {
		return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), nil
	}

	Service.jellyfinFetcher = func(ctx context.Context) (model.CurrentlyPlaying, error) {
		return model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin), nil
	}

	Service.pollInterval = 10 * time.Millisecond

	t.Run("Client Channel Closed Externally", func(t *testing.T) {
		// Clear State
		Service.mu.Lock()
		Service.clients = map[chan string]clientMetadata{}
		Service.mu.Unlock()

		// Start Controller in Background
		done := make(chan struct{})

		go func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/sse", nil)

			Controller.getCurrentPlayingSSE(w, r)

			close(done)
		}()

		// We poll until the client appears. This is faster and safer than Sleep.
		var targetChan chan string

		for {
			Service.mu.RLock()

			for ch := range Service.clients {
				targetChan = ch
				break
			}

			Service.mu.RUnlock()

			if targetChan != nil {
				break // Found it!
			}

			select {
			case <-time.After(1 * time.Second):
				t.Fatal("timed out waiting for client to register")
			case <-time.After(5 * time.Millisecond):
				// Retry per 5ms until we find the channel or timeout
			}
		}

		// Trigger !ok path
		// We manually close the channel. This causes the controller loop to receive (!ok) and exit.
		Service.RemoveClient(context.Background(), targetChan)

		// Assert Exit
		select {
		case <-done:
			// Success: Controller exited cleanly
		case <-time.After(1 * time.Second):
			t.Fatal("handler did not exit after channel closure")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
		r := httptest.NewRequest(http.MethodGet, "/sse", nil)

		done := make(chan bool)

		go func() {
			Controller.getCurrentPlayingSSE(w, r)
			close(done)
		}()

		// Wait for the inevitable write error to trigger 'return'
		select {
		case <-done:
			// Success: Controller exited due to write error
		case <-time.After(1 * time.Second):
			t.Fatal("handler did not exit on write error")
		}
	})

	t.Run("Flush Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/sse", nil)

		// Use a writer that does NOT implement http.Flusher interface.
		// http.NewResponseController(w).Flush() will return ErrNotSupported.
		w := httptest.NewRecorder()
		wrappedW := &test.NonFlusherResponseWriter{ResponseWriter: w}

		ctx, cancel := context.WithCancel(r.Context())

		r = r.WithContext(ctx)

		done := make(chan bool)

		go func() {
			Controller.getCurrentPlayingSSE(wrappedW, r)
			close(done)
		}()

		// Let it run for a cycle to hit the Flush() call
		time.Sleep(10 * time.Millisecond)

		cancel()

		select {
		case <-done:
			// Success
		case <-time.After(1 * time.Second):
			t.Fatal("handler hung")
		}
	})
}
