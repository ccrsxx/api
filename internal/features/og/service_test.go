package og

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/test"
)

func TestService_getOg(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("got %s, want GET", r.Method)
			}

			if r.URL.Query().Get("title") != "hello" {
				t.Errorf("got %s, want query param title=hello", r.URL.Query().Encode())
			}

			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte("image-data"))

			if err != nil {
				t.Fatalf("failed to write response body: %v", err)
			}
		}))

		defer mockServer.Close()

		// Inject the mock server directly
		svc := NewService(ServiceConfig{
			OgUrl:      mockServer.URL,
			HttpClient: mockServer.Client(),
		})

		stream, err := svc.getOg(context.Background(), "title=hello")

		if err != nil {
			t.Fatalf("got error: %v, want success", err)
		}

		defer func() {
			if err := stream.Close(); err != nil {
				t.Errorf("failed to close stream: %v", err)
			}
		}()

		data, err := io.ReadAll(stream)

		if err != nil {
			t.Fatalf("failed to read stream: %v", err)
		}

		if string(data) != "image-data" {
			t.Errorf("got %s, want image-data", string(data))
		}
	})

	t.Run("Status Error (500)", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		svc := NewService(ServiceConfig{
			OgUrl:      mockServer.URL,
			HttpClient: mockServer.Client(),
		})

		_, err := svc.getOg(context.Background(), "")

		if err == nil {
			t.Error("want error for status 500")
		}

		if !strings.Contains(err.Error(), "og request status error: 500") {
			t.Errorf("got unwanted error message: %v", err)
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{
			OgUrl:      "http://127.0.0.1:0", // Invalid port
			HttpClient: &http.Client{Timeout: 1 * time.Millisecond},
		})

		_, err := svc.getOg(context.Background(), "")

		if err == nil {
			t.Error("want network error")
		}

		if !strings.Contains(err.Error(), "og request call error") {
			t.Errorf("got unwanted error message: %v", err)
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{
			OgUrl: "http://\x7f",
		})

		_, err := svc.getOg(context.Background(), "")

		if err == nil {
			t.Error("want error from request creation")
		}
	})

	t.Run("Status Error Body Close Failure", func(t *testing.T) {
		svc := NewService(ServiceConfig{
			OgUrl: "http://example.com",
			HttpClient: &http.Client{
				Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("error")},
					}, nil
				}),
			},
		})

		_, err := svc.getOg(context.Background(), "")

		if err == nil {
			t.Error("want error from body close failure")
		}
	})
}
