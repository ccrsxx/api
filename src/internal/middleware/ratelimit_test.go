package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccrsxx/api/src/internal/api"
	"golang.org/x/time/rate"
)

func TestRateLimit(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Allows requests within limit", func(t *testing.T) {
		mw := RateLimit(5, time.Second)
		server := mw(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "1.1.1.1:1234"

		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", w.Code)
		}
	})

	t.Run("Blocks requests exceeding limit", func(t *testing.T) {
		mw := RateLimit(1, time.Minute)
		server := mw(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "2.2.2.2:1234"

		// 1st request pass
		w1 := httptest.NewRecorder()

		server.ServeHTTP(w1, r)

		// 2nd request block
		w2 := httptest.NewRecorder()

		server.ServeHTTP(w2, r)

		if w2.Code != http.StatusTooManyRequests {
			t.Errorf("second request should fail")
		}

		var res api.ErrorResponse

		if err := json.Unmarshal(w2.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to parse error response: %v", err)
		}

		if res.Error.Message == "" {
			t.Error("want error message")
		}
	})

	t.Run("Headers are set correctly", func(t *testing.T) {
		mw := RateLimit(10, 60*time.Second)

		server := mw(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.RemoteAddr = "127.0.0.1:5000"

		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)

		if w.Header().Get("RateLimit-Limit") != "10" {
			t.Error("wrong limit header")
		}
	})
}

// TestRateLimiter_PruneVisitors tests the logic synchronously
func TestRateLimiter_PruneVisitors(t *testing.T) {
	rl := &rateLimiter{
		visitors: map[string]*visitor{},
	}

	staleIP := "192.168.1.100"

	rl.visitors[staleIP] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now().Add(-2 * time.Minute),
	}

	freshIP := "192.168.1.101"

	rl.visitors[freshIP] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now(),
	}

	// Prune anything older than 1 minute
	rl.pruneVisitors(time.Minute)

	if _, exists := rl.visitors[staleIP]; exists {
		t.Errorf("want stale IP %s to be cleaned up", staleIP)
	}

	if _, exists := rl.visitors[freshIP]; !exists {
		t.Errorf("want fresh IP %s to remain", freshIP)
	}
}

// TestRateLimiter_CleanupLoop tests the goroutine lifecycle
func TestRateLimiter_CleanupLoop(t *testing.T) {
	rl := &rateLimiter{
		visitors: map[string]*visitor{},
	}

	// use separate channels to signal completion and stopping
	// to make sure we can detect if the loop exits properly after receiving the stop signal
	done := make(chan struct{})

	stop := make(chan struct{})

	// Start the cleanup loop with a tiny interval (1ms)
	go func() {
		rl.cleanup(time.Millisecond, stop)

		// Signal that the loop has exited cleanly
		close(done)
	}()

	// Let it run for a bit to ensure it hits the ticker case
	time.Sleep(10 * time.Millisecond)

	close(stop)

	select {
	case <-done:
		// Success will exit here if the loop exits cleanly after receiving the stop signal
	case <-time.After(100 * time.Millisecond):
		t.Fatal("cleanup loop failed to exit after stop signal")
	}
}
