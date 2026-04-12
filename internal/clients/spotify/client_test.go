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

func TestClient_GetCurrentlyPlaying_TokenErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("Token Request Creation Error", func(t *testing.T) {
		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      "http://bad\x7f",
			APIURL:       "http://api",
		})

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("want error from new request creation")
		}
	})

	t.Run("Token Network Error", func(t *testing.T) {
		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      "http://127.0.0.1:0",
			APIURL:       "http://api",
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("want network error")
		}
	})

	t.Run("Token Status 401", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer s.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      s.URL,
			APIURL:       "http://api",
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("got nil, want error from 401 status")
		}
	})

	t.Run("Token Malformed JSON", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer s.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      s.URL,
			APIURL:       "http://api",
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("want error from malformed JSON")
		}
	})

	t.Run("Token Body Close Error", func(t *testing.T) {
		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      "http://auth",
			APIURL:       "http://api",
		})

		c.httpClient.Transport = test.CustomTransport(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &test.ErrorBodyCloser{Reader: strings.NewReader(`{}`)},
				Header:     make(http.Header),
			}, nil
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("want error from token body close")
		}
	})
}

func TestClient_GetCurrentlyPlaying_APIErrors(t *testing.T) {
	ctx := context.Background()

	validAuthHandler := func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"access_token":"valid","expires_in":3600}`)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}

	t.Run("API Request Creation Error", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       "http://bad\x7f",
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("want error from new request creation")
		}
	})

	t.Run("API Network Error", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       "http://",
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("want network error")
		}
	})

	t.Run("API Status 500", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       apiSrv.URL,
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
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       apiSrv.URL,
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("want error from malformed JSON")
		}
	})

	t.Run("API Body Close Error", func(t *testing.T) {
		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      "http://auth",
			APIURL:       "http://api",
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
				Body:       &test.ErrorBodyCloser{Reader: strings.NewReader(`{}`)},
			}, nil
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("want error from API body close")
		}
	})
}

func TestClient_GetCurrentlyPlaying_Logic(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid Item Type", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{"access_token":"v","expires_in":3600}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`{"item": {"name": "JRE", "type": "podcast"}}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       apiSrv.URL,
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
		authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{"access_token":"v","expires_in":3600}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`{"is_playing":true, "item": {"name": "Song", "type": "track"}}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer apiSrv.Close()

		ctx := t.Context()

		memoryCache := cache.NewMemoryCache(ctx, cache.DefaultCleanupInterval)

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       apiSrv.URL,
			MemoryCache:  memoryCache,
		})

		res, err := c.GetCurrentlyPlaying(ctx)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if res.Item.Name != "Song" {
			t.Error("failed to parse response correctly")
		}
	})

	t.Run("Success NoContent", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{"access_token":"v","expires_in":3600}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer apiSrv.Close()

		c := NewClient(Config{
			ClientID:     "id",
			ClientSecret: "sec",
			RefreshToken: "ref",
			AuthURL:      authSrv.URL,
			APIURL:       apiSrv.URL,
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if !errors.Is(err, ErrNoContent) {
			t.Fatalf("got %v, want ErrNoContent", err)
		}
	})
}
