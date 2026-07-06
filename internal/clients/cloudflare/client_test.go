package cloudflare_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/cloudflare"
	"github.com/ccrsxx/api/internal/test"
)

func TestNewClient(t *testing.T) {
	client := cloudflare.NewClient(cloudflare.Config{})

	if client == nil {
		t.Fatal("want client to be initialized, got nil")
	}
}

func TestClient_VerifyTurnstile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{"success": true}`))

			if err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := cloudflare.NewClient(cloudflare.Config{APIURL: mockServer.URL})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}
	})

	t.Run("Marshal Error", func(t *testing.T) {
		restore := cloudflare.SetJSONMarshal(func(v any) ([]byte, error) {
			return nil, errors.New("marshal error")
		})

		defer restore()

		c := cloudflare.NewClient(cloudflare.Config{APIURL: "http://example.com"})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err == nil {
			t.Error("want marshal error")
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := cloudflare.NewClient(cloudflare.Config{APIURL: "http://bad\x7f"})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err == nil {
			t.Error("want error from NewRequestWithContext")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := cloudflare.NewClient(cloudflare.Config{APIURL: "http://invalid.url.local"})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err == nil {
			t.Error("want error from httpClient.Do")
		}
	})

	t.Run("Status 500", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		c := cloudflare.NewClient(cloudflare.Config{APIURL: mockServer.URL})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err == nil {
			t.Error("want status error for 500")
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

		c := cloudflare.NewClient(cloudflare.Config{APIURL: mockServer.URL})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err == nil {
			t.Error("want decode error")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := cloudflare.NewClient(cloudflare.Config{
			APIURL: "http://example.com",
			HTTPClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader(`{"success": true}`)},
						Header:     make(http.Header),
					}, nil
				}),
			},
		})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err != nil {
			t.Fatalf("got: %v, want VerifyTurnstile to handle body close error gracefully", err)
		}
	})

	t.Run("Captcha Failure", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{"success": false}`))

			if err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := cloudflare.NewClient(cloudflare.Config{APIURL: mockServer.URL})

		err := c.VerifyTurnstile(context.Background(), "test_token", "127.0.0.1")

		if err == nil {
			t.Fatal("want error for failed captcha")
		}

		httpErr, ok := err.(*api.HTTPError)

		if !ok {
			t.Fatalf("want *api.HTTPError, got %T", err)
		}

		if httpErr.StatusCode != http.StatusForbidden {
			t.Errorf("got status %d, want %d", httpErr.StatusCode, http.StatusForbidden)
		}
	})
}
