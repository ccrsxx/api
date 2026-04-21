package spotify

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/cache"
	"github.com/ccrsxx/api/internal/test"
)

func TestNewClient(t *testing.T) {
	client := NewClient(Config{})

	if client == nil {
		t.Fatal("want client to be initialized")
	}
}

func TestClient_GetCurrentlyPlaying(t *testing.T) {
	ctx := context.Background()

	validAuthHandler := func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"access_token":"valid","expires_in":3600}`)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}

	t.Run("Token Error Bubbles Up", func(t *testing.T) {
		c := NewClient(Config{AuthURL: "http://bad\x7f"})

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("want token fetch error to bubble up")
		}
	})

	t.Run("API Request Creation Error", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		c := NewClient(Config{
			AuthURL: authSrv.URL,
			APIURL:  "http://bad\x7f",
		})

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("want error from new request creation")
		}
	})

	t.Run("API Network Error", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		c := NewClient(Config{
			AuthURL: authSrv.URL,
			APIURL:  "http://invalid.url.local",
		})

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("want network error")
		}
	})

	t.Run("API Body Close Error", func(t *testing.T) {
		c := NewClient(Config{
			AuthURL: "http://auth",
			APIURL:  "http://api",
		})

		c.httpClient.Transport = test.CustomTransport(func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "auth") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"t","expires_in":3600}`)),
				}, nil
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("[]")},
			}, nil
		})

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Fatal("want error from API body close")
		}
	})

	t.Run("Success NoContent", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			AuthURL: authSrv.URL,
			APIURL:  apiSrv.URL,
		})

		if _, err := c.GetCurrentlyPlaying(ctx); !errors.Is(err, ErrNoContent) {
			t.Fatalf("got %v, want ErrNoContent", err)
		}
	})

	t.Run("API Status Error", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			AuthURL: authSrv.URL,
			APIURL:  apiSrv.URL,
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("got nil, want 500 error")
		}

		if !strings.Contains(err.Error(), "status error: 500") {
			t.Errorf("got %v, want status error message", err)
		}
	})

	t.Run("API Malformed JSON", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			AuthURL: authSrv.URL,
			APIURL:  apiSrv.URL,
		})

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("want error from malformed JSON")
		}
	})

	t.Run("Invalid Item Type", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"item": {"name": "JRE", "type": "podcast"}}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			AuthURL: authSrv.URL,
			APIURL:  apiSrv.URL,
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("got nil, want invalid item type error")
		}

		if !strings.Contains(err.Error(), "invalid item type") {
			t.Errorf("wrong error message: %v", err)
		}
	})

	t.Run("Success Track", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"is_playing":true, "item": {"name": "Song", "type": "track"}}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			AuthURL:     authSrv.URL,
			APIURL:      apiSrv.URL,
			MemoryCache: cache.NewMemoryCache(ctx, cache.DefaultCleanupInterval),
		})

		res, err := c.GetCurrentlyPlaying(ctx)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if res.Item.Name != "Song" {
			t.Error("failed to parse response correctly")
		}
	})
}

func TestClient_GetAccessToken(t *testing.T) {
	ctx := context.Background()

	t.Run("Request Creation Error", func(t *testing.T) {
		c := NewClient(Config{AuthURL: "http://bad\x7f"})

		if _, err := c.getAccessToken(ctx); err == nil {
			t.Error("want error from new request creation")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := NewClient(Config{AuthURL: "http://127.0.0.1:0"})

		if _, err := c.getAccessToken(ctx); err == nil {
			t.Error("want network error")
		}
	})

	t.Run("Status Error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer s.Close()

		c := NewClient(Config{AuthURL: s.URL})

		if _, err := c.getAccessToken(ctx); err == nil {
			t.Fatal("want error from non-200 status")
		}
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer s.Close()

		c := NewClient(Config{AuthURL: s.URL})

		if _, err := c.getAccessToken(ctx); err == nil {
			t.Error("want error from malformed JSON")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := NewClient(Config{AuthURL: "http://auth"})

		c.httpClient.Transport = test.CustomTransport(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("[]")},
				Header:     make(http.Header),
			}, nil
		})

		if _, err := c.getAccessToken(ctx); err == nil {
			t.Fatal("want error bubbled or logged from token body close")
		}
	})

	t.Run("Success", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{"access_token":"valid_token","expires_in":3600}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer s.Close()

		c := NewClient(Config{AuthURL: s.URL})

		token, err := c.getAccessToken(ctx)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}
		if token != "valid_token" {
			t.Errorf("got %q, want valid_token", token)
		}
	})
}
