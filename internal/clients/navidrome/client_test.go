package navidrome_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/navidrome"
	"github.com/ccrsxx/api/internal/test"
)

func TestNewClient(t *testing.T) {
	client := navidrome.NewClient(navidrome.Config{})

	if client == nil {
		t.Fatal("want client to be initialized, got nil")
	}
}

func TestClient_GetNowPlaying(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			resp := navidrome.JSONWrapper{
				Subsonic: navidrome.Subsonic{
					NowPlaying: &navidrome.NowPlaying{
						Entry: []navidrome.NowPlayingEntry{
							{UserName: "user1"},
						},
					},
				},
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		entries, err := c.GetNowPlaying(context.Background())

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if len(entries) != 1 {
			t.Errorf("got %d, want 1 entry", len(entries))
		}

		if entries[0].UserName != "user1" {
			t.Errorf("got %q, want %q", entries[0].UserName, "user1")
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := navidrome.NewClient(navidrome.Config{URL: "http://bad\x7f"})

		_, err := c.GetNowPlaying(context.Background())

		if err == nil {
			t.Error("want error from NewRequestWithContext")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := navidrome.NewClient(navidrome.Config{URL: "http://invalid.url.local"})

		_, err := c.GetNowPlaying(context.Background())

		if err == nil {
			t.Error("want error from httpClient.Do")
		}
	})

	t.Run("Status Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetNowPlaying(context.Background())

		if err == nil {
			t.Error("want status error for 401")
		}
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetNowPlaying(context.Background())

		if err == nil {
			t.Error("want decode error")
		}
	})

	t.Run("Nil NowPlaying", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			resp := navidrome.JSONWrapper{
				Subsonic: navidrome.Subsonic{},
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetNowPlaying(context.Background())

		if err == nil {
			t.Error("want error for nil now playing")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		resp := navidrome.JSONWrapper{
			Subsonic: navidrome.Subsonic{
				NowPlaying: &navidrome.NowPlaying{
					Entry: []navidrome.NowPlayingEntry{},
				},
			},
		}

		body, _ := json.Marshal(resp)

		c := navidrome.NewClient(navidrome.Config{
			URL: "http://example.com",
			HTTPClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader(string(body))},
						Header:     make(http.Header),
					}, nil
				}),
			},
		})

		_, err := c.GetNowPlaying(context.Background())

		if err != nil {
			t.Fatalf("got: %v, want GetNowPlaying to handle body close error gracefully", err)
		}
	})
}

func TestClient_GetCoverArtStream(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/webp")
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte("fake-image-bytes")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		body, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		defer func() {
			if err := body.Close(); err != nil {
				t.Errorf("failed to close body: %v", err)
			}
		}()

		data, err := io.ReadAll(body)

		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if string(data) != "fake-image-bytes" {
			t.Errorf("got %q, want %q", string(data), "fake-image-bytes")
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := navidrome.NewClient(navidrome.Config{URL: "http://bad\x7f"})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Error("want error from NewRequestWithContext")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := navidrome.NewClient(navidrome.Config{URL: "http://invalid.url.local"})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Error("want error from httpClient.Do")
		}
	})

	t.Run("Status Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Error("want status error for 500")
		}
	})

	t.Run("Status Error Body Close Error", func(t *testing.T) {
		c := navidrome.NewClient(navidrome.Config{
			URL: "http://example.com",
			HTTPClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("")},
						Header:     make(http.Header),
					}, nil
				}),
			},
		})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Error("want status error even when body close fails")
		}
	})

	t.Run("JSON Content Type Cover Art Not Found", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := navidrome.JSONWrapper{
				Subsonic: navidrome.Subsonic{
					Error: &navidrome.Error{Code: 70, Message: "data not found"},
				},
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Fatal("want error for cover art not found")
		}

		var httpErr *api.HTTPError

		if !errors.As(err, &httpErr) {
			t.Fatalf("want *api.HTTPError, got %T", err)
		}

		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("got status %d, want %d", httpErr.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("JSON Content Type Other Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := navidrome.JSONWrapper{
				Subsonic: navidrome.Subsonic{
					Error: &navidrome.Error{Code: 10, Message: "required parameter is missing"},
				},
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Error("want error for non-70 error code")
		}

		var httpErr *api.HTTPError

		if errors.As(err, &httpErr) {
			t.Error("want generic error, not *api.HTTPError")
		}
	})

	t.Run("JSON Content Type Malformed JSON", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := navidrome.NewClient(navidrome.Config{URL: mockServer.URL})

		_, err := c.GetCoverArtStream(context.Background(), "cover-1")

		if err == nil {
			t.Error("want decode error")
		}
	})
}
