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
	setupTest := func(ctx context.Context) (*Controller, *Service) {
		dummySpotify := &mockDataFetcher{
			result: model.NewDefaultCurrentlyPlaying(model.PlatformSpotify),
		}

		dummyJellyfin := &mockDataFetcher{
			result: model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin),
		}

		svc := NewService(ServiceConfig{
			PollInterval:    10 * time.Millisecond,
			SpotifyService:  dummySpotify,
			JellyfinService: dummyJellyfin,
		})

		ctrl := NewController(ctx, svc)

		return ctrl, svc
	}

	t.Run("Client Channel Closed Externally", func(t *testing.T) {
		ctx := t.Context()

		ctrl, svc := setupTest(ctx)

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
		ctx := t.Context()

		ctrl, _ := setupTest(ctx)

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
		ctx := t.Context()

		ctrl, _ := setupTest(ctx)

		r := httptest.NewRequest(http.MethodGet, "/sse", nil)

		// Use a writer that does NOT implement http.Flusher interface.
		// http.NewResponseController(w).Flush() will return ErrNotSupported.
		w := &test.NonFlusherResponseWriter{ResponseWriter: httptest.NewRecorder()}

		ctx, cancel := context.WithCancel(r.Context())

		defer cancel()

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

	t.Run("App Shutdown Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())

		defer cancel()

		ctrl, _ := setupTest(ctx)

		req := httptest.NewRequest(http.MethodGet, "/sse", nil)
		rec := httptest.NewRecorder()

		done := make(chan struct{})

		go func() {
			ctrl.getCurrentPlayingSSE(rec, req)
			close(done) // This channel closes ONLY when the handler successfully returns
		}()

		time.Sleep(10 * time.Millisecond)

		// HIT THE KILL SWITCH! This triggers the <-appShutdown case
		cancel()

		select {
		case <-done:
			// Success! The handler caught the shutdown and exited cleanly.
		case <-time.After(1 * time.Second):
			// If we get here, your server is hanging!
			t.Fatal("handler did not exit after app context was canceled (goroutine leak!)")
		}
	},
	)
}
