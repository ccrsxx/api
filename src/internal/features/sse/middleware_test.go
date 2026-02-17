package sse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	handler := Middleware.IsConnectionAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("Connection Allowed", func(t *testing.T) {
		// Lock required for test stability: The background worker may still be
		// reading these configuration fields while we restore them.

		Service.mu.Lock()

		Service.clients = map[chan string]clientMetadata{}
		Service.ipAddressCounts = map[string]int{}

		Service.mu.Unlock()

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.RemoteAddr = "1.1.1.1:1234"

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200 ok", rec.Code)
		}

		wantHeaders := map[string]string{
			"Connection":    "keep-alive",
			"Content-Type":  "text/event-stream",
			"Cache-Control": "no-cache",
		}

		for k, v := range wantHeaders {
			if got := rec.Header().Get(k); got != v {
				t.Errorf("got header %s: %s, want %s", k, v, got)
			}
		}
	})

	t.Run("IP Limit Reached", func(t *testing.T) {
		targetIP := "2.2.2.2"

		Service.mu.Lock()

		// Max out the IP limit for the target IP
		Service.ipAddressCounts[targetIP] = maxClientsPerIP

		Service.mu.Unlock()

		defer func() {
			Service.mu.Lock()

			delete(Service.ipAddressCounts, targetIP)

			Service.mu.Unlock()
		}()

		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.RemoteAddr = targetIP + ":12345"

		handler.ServeHTTP(w, r)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("got %d, want 429 too many requests", w.Code)
		}
	})
}
