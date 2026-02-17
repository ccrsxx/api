package og

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestService_getOg(t *testing.T) {
	originalClient := Service.httpClient

	originalDev := config.Config().IsDevelopment

	defer func() {
		Service.httpClient = originalClient
		config.Config().IsDevelopment = originalDev
	}()

	t.Run("Success Production URL", func(t *testing.T) {
		config.Config().IsDevelopment = false

		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
				// Verify Production URL
				if !strings.HasPrefix(req.URL.String(), "http://10.0.0.60:4444/og") {
					t.Errorf("want prod url, got %s", req.URL.String())
				}

				if req.URL.Query().Get("title") != "hello" {
					t.Errorf("want query param, got %s", req.URL.Query().Encode())
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("image-data")),
					Header:     make(http.Header),
				}, nil
			}),
		}

		stream, err := Service.getOg(context.Background(), "title=hello")

		if err != nil {
			t.Fatalf("unwant error: %v", err)
		}

		defer func() {
			if err := stream.Close(); err != nil {
				t.Errorf("unwant error closing stream: %v", err)
			}
		}()

		data, _ := io.ReadAll(stream)

		if string(data) != "image-data" {
			t.Errorf("want image-data, got %s", data)
		}

	})

	t.Run("Success Development URL", func(t *testing.T) {
		config.Config().IsDevelopment = true

		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
				if !strings.HasPrefix(req.URL.String(), "http://localhost:4444/og") {
					t.Errorf("want dev url, got %s", req.URL.String())
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("dev-image")),
					Header:     make(http.Header),
				}, nil
			}),
		}

		_, err := Service.getOg(context.Background(), "")

		if err != nil {
			t.Fatalf("unwant error: %v", err)
		}
	})

	t.Run("Request Creation Error", func(t *testing.T) {
		// Simulate request creation error by passing an invalid URL
		_, err := Service.getOg(nil, "")

		if err == nil {
			t.Error("want error from nil context")
		}
	})

	t.Run("Network Call Error", func(t *testing.T) {
		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network down")
			}),
		}

		_, err := Service.getOg(context.Background(), "")

		if err == nil {
			t.Error("want network error")
		}

		if !strings.Contains(err.Error(), "og request call error") {
			t.Errorf("wrong error message: %v", err)
		}
	})

	t.Run("Status Error (500)", func(t *testing.T) {
		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("error")),
				}, nil
			}),
		}

		_, err := Service.getOg(context.Background(), "")

		if err == nil {
			t.Error("want status error")
		}

		if !strings.Contains(err.Error(), "og request status error: 500") {
			t.Errorf("wrong error message: %v", err)
		}
	})

	t.Run("Status Error Body Close Failure", func(t *testing.T) {
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
			t.Error("want status error")
		}
	})
}
