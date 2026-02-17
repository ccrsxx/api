package og

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestService_getOg(t *testing.T) {
	originalOgUrl := Service.ogUrl
	originalClient := Service.httpClient

	originalDev := config.Config().IsDevelopment

	defer func() {
		Service.ogUrl = originalOgUrl
		Service.httpClient = originalClient

		config.Config().IsDevelopment = originalDev
	}()

	// Default to Production for most tests so we can inject the mock URL
	config.Config().IsDevelopment = false

	t.Run("Success (Production URL)", func(t *testing.T) {
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

		Service.ogUrl = mockServer.URL
		Service.httpClient = mockServer.Client()

		stream, err := Service.getOg(context.Background(), "title=hello")

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

	t.Run("Success (Development URL)", func(t *testing.T) {
		// Toggle Dev Mode ON just for this sub-test
		config.Config().IsDevelopment = true

		defer func() {
			config.Config().IsDevelopment = false
		}()

		// We MUST use CustomTransport here.
		// Why? Because in Dev mode, the code hardcodes "localhost:4444".
		// We can't make the code call a dynamic httptest port, so we intercept the client instead.

		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
				// Assert the code actually tried to hit the hardcoded dev URL
				expectedPrefix := "http://localhost:4444/og"

				if !strings.HasPrefix(req.URL.String(), expectedPrefix) {
					t.Errorf("got %s, want dev url prefix %s", req.URL.String(), expectedPrefix)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("dev-image-data")),
					Header:     make(http.Header),
				}, nil
			}),
		}

		stream, err := Service.getOg(context.Background(), "title=dev")

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

		if string(data) != "dev-image-data" {
			t.Errorf("got %s, want dev-image-data", string(data))
		}
	})

	t.Run("Status Error (500)", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		Service.ogUrl = mockServer.URL
		Service.httpClient = mockServer.Client()

		_, err := Service.getOg(context.Background(), "")

		if err == nil {
			t.Error("want error for status 500")
		}

		if !strings.Contains(err.Error(), "og request status error: 500") {
			t.Errorf("got unexpected error message: %v", err)
		}
	})

	t.Run("Network Error", func(t *testing.T) {
		Service.ogUrl = "http://127.0.0.1:0" // Invalid port
		Service.httpClient = &http.Client{Timeout: 1 * time.Millisecond}

		_, err := Service.getOg(context.Background(), "")

		if err == nil {
			t.Error("want network error")
		}

		if !strings.Contains(err.Error(), "og request call error") {
			t.Errorf("got unexpected error message: %v", err)
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		Service.ogUrl = "http://\x7f"

		_, err := Service.getOg(context.Background(), "")

		if err == nil {
			t.Error("want error from request creation")
		}
	})

	t.Run("Status Error Body Close Failure", func(t *testing.T) {
		Service.ogUrl = "http://example.com"

		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("error")},
				}, nil
			}),
		}

		_, err := Service.getOg(context.Background(), "")

		if err == nil {
			t.Error("want error from body close failure")
		}
	})
}
