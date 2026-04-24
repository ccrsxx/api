package jellyfin_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/clients/jellyfin"
	"github.com/ccrsxx/api/internal/test"
)

func TestNewClient(t *testing.T) {
	client := jellyfin.NewClient(jellyfin.Config{})

	if client == nil {
		t.Fatal("want client to be initialized, got nil")
	}
}

func TestClient_GetSessions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			id := "1"

			if err := json.NewEncoder(w).Encode([]jellyfin.SessionInfo{{ID: &id}}); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := jellyfin.NewClient(jellyfin.Config{URL: mockServer.URL})

		sessions, err := c.GetSessions(context.Background())

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if len(sessions) != 1 {
			t.Errorf("got %d, want 1 session", len(sessions))
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := jellyfin.NewClient(jellyfin.Config{URL: "http://bad\x7f"})

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("want error from NewRequestWithContext")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := jellyfin.NewClient(jellyfin.Config{URL: "http://invalid.url.local"})

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

		c := jellyfin.NewClient(jellyfin.Config{URL: mockServer.URL})

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

		c := jellyfin.NewClient(jellyfin.Config{URL: mockServer.URL})

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("want decode error")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := jellyfin.NewClient(jellyfin.Config{
			URL: "http://example.com",
			HTTPClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("[]")},
						Header:     make(http.Header),
					}, nil
				}),
			},
		})

		_, err := c.GetSessions(context.Background())

		if err != nil {
			t.Fatalf("got: %v, want GetSessions to handle body close error gracefully", err)
		}
	})
}
