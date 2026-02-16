package jellyfin

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Helper for pointer creation
func stringPtr(s string) *string {
	return &s
}

func TestDefaultClient(t *testing.T) {
	c1 := DefaultClient()

	if c1 == nil {
		t.Fatal("expected default client, got nil")
	}

	c2 := DefaultClient()
	if c1 != c2 {
		t.Error("DefaultClient should return the same singleton instance")
	}
}

func TestClient_GetSessions(t *testing.T) {
	t.Run("Success Path", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode([]SessionInfo{{Id: stringPtr("1")}}); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := New(mockServer.URL, "key", "img", "user")

		sessions, err := c.GetSessions(context.Background())

		if err != nil || len(sessions) != 1 {
			t.Fatalf("expected success, got err: %v", err)
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		// A URL with a control character triggers a NewRequest error
		c := New("http://localhost\x7f", "key", "img", "user")

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("expected error from NewRequestWithContext")
		}
	})

	t.Run("HTTP Do Error", func(t *testing.T) {
		c := New("http://invalid.url.local", "key", "img", "user")

		// Short timeout to speed up the failure
		c.httpClient.Timeout = 10 * time.Millisecond

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("expected error from httpClient.Do")
		}
	})

	t.Run("Status Not OK Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer mockServer.Close()

		c := New(mockServer.URL, "key", "img", "user")

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("expected status error for 401")
		}
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid-json`))
		}))

		defer mockServer.Close()

		c := New(mockServer.URL, "key", "img", "user")

		_, err := c.GetSessions(context.Background())

		if err == nil {
			t.Error("expected decode error")
		}
	})
}

// Fixed Body Close Error Test to prevent hanging
type errorCloser struct {
	io.Reader
}

func (e *errorCloser) Close() error {
	return errors.New("forced close error")
}

func TestClient_GetSessions_CloseError(t *testing.T) {
	c := New("http://localhost", "key", "img", "user")

	// Inject a response with a body that fails on Close()
	// We provide an empty JSON array so the decoder finishes immediately and calls Close()
	c.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &errorCloser{Reader: strings.NewReader("[]")},
			Header:     make(http.Header),
		}, nil
	})

	// This triggers the defer func() { if err := res.Body.Close() ... }() block
	_, err := c.GetSessions(context.Background())

	if err != nil {
		t.Fatalf("expected GetSessions to handle body close error gracefully, got: %v", err)
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
