package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/api"
	"golang.org/x/time/rate"
)

func TestRateLimit(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Allows requests within limit", func(t *testing.T) {
		ctx := t.Context()

		mw := RateLimit(ctx, 5, time.Second)
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
		ctx := t.Context()

		mw := RateLimit(ctx, 1, time.Minute)
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
			t.Fatalf("second request should fail")
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
		ctx := t.Context()

		mw := RateLimit(ctx, 10, 60*time.Second)

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

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := &rateLimiter{
		visitors: map[string]*visitor{},
	}

	staleIP := "192.168.1.100"

	rl.visitors[staleIP] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now().Add(-2 * time.Minute), // Way in the past, gets cleaned up
	}

	freshIP := "192.168.1.101"

	rl.visitors[freshIP] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now().Add(4 * time.Minute), // In the future, should not be cleaned up
	}

	ctx := t.Context()

	// Start cleanup with a 10ms interval
	go rl.cleanup(ctx, 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond) // Wait for cleanup to run

	// Lock required to prevent race condition
	// We make sure that cleanup function has unlocked before we check the map
	// otherwise we might be checking while it's still locked and gets a incorrect result
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, exists := rl.visitors[staleIP]; exists {
		t.Errorf("want stale IP %s to be cleaned up", staleIP)
	}

	if _, exists := rl.visitors[freshIP]; !exists {
		t.Errorf("want fresh IP %s to remain", freshIP)
	}
}
