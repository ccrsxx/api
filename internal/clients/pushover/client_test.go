package pushover_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/clients/pushover"
	"github.com/ccrsxx/api/internal/test"
)

func TestNewClient(t *testing.T) {
	client := pushover.NewClient(pushover.Config{})

	if client == nil {
		t.Fatal("want client to be initialized, got nil")
	}
}

func TestClient_SendMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{"status": 1, "request": "test-request-id"}`))

			if err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))

		defer mockServer.Close()

		c := pushover.NewClient(pushover.Config{APIURL: mockServer.URL})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{
			Message: "test message",
			Title:   "test title",
		})

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}
	})

	t.Run("Marshal Error", func(t *testing.T) {
		restore := pushover.SetJSONMarshal(func(v any) ([]byte, error) {
			return nil, errors.New("marshal error")
		})

		defer restore()

		c := pushover.NewClient(pushover.Config{APIURL: "http://example.com"})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{})

		if err == nil {
			t.Error("want marshal error")
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		c := pushover.NewClient(pushover.Config{APIURL: "http://bad\x7f"})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{})

		if err == nil {
			t.Error("want error from NewRequestWithContext")
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		c := pushover.NewClient(pushover.Config{APIURL: "http://invalid.url.local"})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{})

		if err == nil {
			t.Error("want error from httpClient.Do")
		}
	})

	t.Run("Status 500", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		c := pushover.NewClient(pushover.Config{APIURL: mockServer.URL})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{})

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

		c := pushover.NewClient(pushover.Config{APIURL: mockServer.URL})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{})

		if err == nil {
			t.Error("want decode error")
		}
	})

	t.Run("Body Close Error", func(t *testing.T) {
		c := pushover.NewClient(pushover.Config{
			APIURL: "http://example.com",
			HTTPClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader(`{"status": 1}`)},
						Header:     make(http.Header),
					}, nil
				}),
			},
		})

		err := c.SendMessage(context.Background(), pushover.MessageRequest{})

		if err != nil {
			t.Fatalf("got: %v, want SendMessage to handle body close error gracefully", err)
		}
	})
}
