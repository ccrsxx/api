package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestNewClient(t *testing.T) {
	client := NewClient(Config{})

	if client == nil {
		t.Error("want client to be initialized")
	}
}

func TestClient_GetCurrentUser(t *testing.T) {
	ctx := context.Background()
	token := "dummy_token"

	t.Run("Request Creation Error", func(t *testing.T) {
		c := NewClient(Config{APIURL: "http://bad\x7f"})

		if _, err := c.GetCurrentUser(ctx, token); err == nil {
			t.Error("want error from new request creation")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := NewClient(Config{APIURL: "http://invalid.url.local"})

		if _, err := c.GetCurrentUser(ctx, token); err == nil {
			t.Error("want network error")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := NewClient(Config{APIURL: "http://api"})

		c.httpClient.Transport = test.CustomTransport(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("{}")},
				Header:     make(http.Header),
			}, nil
		})

		if _, err := c.GetCurrentUser(ctx, token); err != nil {
			t.Errorf("want nil error despite body close error, got: %v", err)
		}
	})

	t.Run("Status Error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized) // Any non-200 status
		}))

		defer s.Close()

		c := NewClient(Config{APIURL: s.URL})

		_, err := c.GetCurrentUser(ctx, token)

		if err == nil {
			t.Fatal("want error from non-200 status")
		}

		if !strings.Contains(err.Error(), "status error") {
			t.Errorf("got %v, want status error message", err)
		}
	})

	t.Run("Decode Error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer s.Close()

		c := NewClient(Config{APIURL: s.URL})

		if _, err := c.GetCurrentUser(ctx, token); err == nil {
			t.Error("want error from malformed JSON decode")
		}
	})

	t.Run("Success", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
				t.Fatalf("got Accept %q, want application/vnd.github+json", got)
			}

			if got := r.Header.Get("Authorization"); got != "Bearer "+token {
				t.Fatalf("got Authorization %q, want Bearer %s", got, token)
			}

			if got := r.Header.Get("X-GitHub-Api-Version"); got != defaultGithubAPIVersion {
				t.Fatalf("got API Version %q, want %s", got, defaultGithubAPIVersion)
			}

			if _, err := w.Write([]byte("{}")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer s.Close()

		c := NewClient(Config{APIURL: s.URL})

		if _, err := c.GetCurrentUser(ctx, token); err != nil {
			t.Fatalf("unwanted error: %v", err)
		}
	})
}
