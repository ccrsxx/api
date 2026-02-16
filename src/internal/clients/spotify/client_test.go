package spotify

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultClient(t *testing.T) {
	c := DefaultClient()

	if c == nil {
		t.Fatal("expected default client to be initialized")
	}
}

func TestClient_GetCurrentlyPlaying_TokenErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("Token Request Creation Error", func(t *testing.T) {
		c := New("id", "sec", "ref", "http://bad\x7f", "http://api")

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("expected error from new request creation")
		}
	})

	t.Run("Token Network Error", func(t *testing.T) {
		c := New("id", "sec", "ref", "http://127.0.0.1:0", "http://api")

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("expected network error")
		}
	})

	t.Run("Token Status 401", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer s.Close()

		c := New("id", "sec", "ref", s.URL, "http://api")

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("expected error from 401 status")
		}
	})

	t.Run("Token Malformed JSON", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(`{bad-json`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer s.Close()

		c := New("id", "sec", "ref", s.URL, "http://api")

		if _, err := c.GetCurrentlyPlaying(ctx); err == nil {
			t.Error("expected json decode error")
		}
	})

	t.Run("Token Body Close Error", func(t *testing.T) {
		c := New("id", "sec", "ref", "http://auth", "http://api")

		c.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &errorCloser{Reader: strings.NewReader(`{}`)},
				Header:     make(http.Header),
			}, nil
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("expected error from token body close")
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

		c := New("id", "sec", "ref", authSrv.URL, "http://bad\x7f")

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("expected error from new request creation")
		}
	})

	t.Run("API Network Error", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		c := New("id", "sec", "ref", authSrv.URL, "http://127.0.0.1:0")

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("expected network error")
		}
	})

	t.Run("API Status 500", func(t *testing.T) {
		authSrv := httptest.NewServer(http.HandlerFunc(validAuthHandler))

		defer authSrv.Close()

		apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer apiSrv.Close()

		c := New("id", "sec", "ref", authSrv.URL, apiSrv.URL)

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("expected 500 error, got nil")
		}

		if !strings.Contains(err.Error(), "status error: 500") {
			t.Errorf("expected status error message, got: %v", err)
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

		c := New("id", "sec", "ref", authSrv.URL, apiSrv.URL)

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Error("expected error from malformed JSON")
		}
	})

	t.Run("API Body Close Error", func(t *testing.T) {
		c := New("id", "sec", "ref", "http://auth", "http://api")

		c.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "auth") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"t","expires_in":3600}`)),
				}, nil
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &errorCloser{Reader: strings.NewReader(`{}`)},
			}, nil
		})

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("expected error from API body close")
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

		c := New("id", "sec", "ref", authSrv.URL, apiSrv.URL)

		_, err := c.GetCurrentlyPlaying(ctx)

		if err == nil {
			t.Fatal("expected invalid item type error, got nil")
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

		c := New("id", "sec", "ref", authSrv.URL, apiSrv.URL)

		resp, err := c.GetCurrentlyPlaying(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp == nil || resp.Item.Name != "Song" {
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

		c := New("id", "sec", "ref", authSrv.URL, apiSrv.URL)

		resp, err := c.GetCurrentlyPlaying(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp != nil {
			t.Error("expected nil response for 204 No Content")
		}
	})
}

type errorCloser struct {
	io.Reader
}

func (e *errorCloser) Close() error {
	return errors.New("forced close error")
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
