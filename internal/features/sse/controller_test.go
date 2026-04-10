package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/test"
)

func TestController_getCurrentPlayingSSE(t *testing.T) {
	setupTest := func() (*Service, *Controller) {
		dummySpotify := func(ctx context.Context) (model.CurrentlyPlaying, error) {
			return model.NewDefaultCurrentlyPlaying(model.PlatformSpotify), nil
		}

		dummyJellyfin := func(ctx context.Context) (model.CurrentlyPlaying, error) {
			return model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin), nil
		}

		svc := NewService(ServiceConfig{
			PollInterval:    10 * time.Millisecond,
			SpotifyFetcher:  dummySpotify,
			JellyfinFetcher: dummyJellyfin,
		})

		ctrl := NewController(svc)

		return svc, ctrl
	}

	t.Run("Client Channel Closed Externally", func(t *testing.T) {
		svc, ctrl := setupTest()

		done := make(chan struct{})

		go func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/sse", nil)

			ctrl.getCurrentPlayingSSE(w, r)

			close(done)
		}()

		// We poll until the client appears. This is faster and safer than Sleep.
		var targetChan chan string

		for {
			svc.mu.RLock()

			for ch := range svc.clients {
				targetChan = ch
				break
			}

			svc.mu.RUnlock()

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
		svc.RemoveClient(context.Background(), targetChan)

		// Assert Exit
		select {
		case <-done:
			// Success: Controller exited cleanly
		case <-time.After(1 * time.Second):
			t.Fatal("handler did not exit after channel closure")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		_, ctrl := setupTest()

		w := &test.ErrorResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
		r := httptest.NewRequest(http.MethodGet, "/sse", nil)

		done := make(chan bool)

		go func() {
			ctrl.getCurrentPlayingSSE(w, r)
			close(done)
		}()

		// Wait for the inevitable write error to trigger 'return'
		select {
		case <-done:
			// Success: Controller exited due to write error
		case <-time.After(1 * time.Second):
			t.Fatal("handler did not exit on write error")
		}

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("Flush Error", func(t *testing.T) {
		_, ctrl := setupTest()

		r := httptest.NewRequest(http.MethodGet, "/sse", nil)

		// Use a writer that does NOT implement http.Flusher interface.
		// http.NewResponseController(w).Flush() will return ErrNotSupported.
		w := &test.NonFlusherResponseWriter{ResponseWriter: httptest.NewRecorder()}

		ctx, cancel := context.WithCancel(r.Context())

		r = r.WithContext(ctx)

		done := make(chan struct{})

		go func() {
			ctrl.getCurrentPlayingSSE(w, r)
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
