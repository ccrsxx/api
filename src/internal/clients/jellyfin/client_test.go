package jellyfin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/src/internal/test"
)

func TestDefaultClient(t *testing.T) {
	client := DefaultClient()

	if client == nil {
		t.Fatal("want default client, got nil")
	}
}

func TestClient_GetSessions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			id := "1"

			if err := json.NewEncoder(w).Encode([]SessionInfo{{Id: &id}}); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := New(mockServer.URL, "key", "img", "user")

		sessions, err := c.GetSessions(context.Background())

		if err != nil {
			t.Fatalf("want success, got err: %v", err)
		}

		if len(sessions) != 1 {
			t.Errorf("want 1 session, got %d", len(sessions))
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := New("http://localhost\x7f", "key", "img", "user")

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("want error from NewRequestWithContext")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := New("http://invalid.url.local", "key", "img", "user")

		c.httpClient.Timeout = 10 * time.Millisecond

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("want error from httpClient.Do")
		}
	})

	t.Run("Status 401", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer mockServer.Close()

		c := New(mockServer.URL, "key", "img", "user")

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("want status error for 401")
		}
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`invalid-json`))

			if err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := New(mockServer.URL, "key", "img", "user")

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("want decode error")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := New("http://localhost", "key", "img", "user")

		c.httpClient.Transport = test.CustomTransport(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("[]")},
				Header:     make(http.Header),
			}, nil
		})

		_, err := c.GetSessions(context.Background())

		if err != nil {
			t.Fatalf("want GetSessions to handle body close error gracefully, got: %v", err)
		}
	})
}
