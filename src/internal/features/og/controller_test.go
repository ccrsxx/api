package og

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestController_getOg(t *testing.T) {
	originalClient := Service.httpClient

	originalProd := config.Config().IsProduction

	defer func() {
		Service.httpClient = originalClient
		config.Config().IsProduction = originalProd
	}()

	t.Run("Success Default (No Cache)", func(t *testing.T) {
		config.Config().IsProduction = false

		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("png-data")),
					Header:     make(http.Header),
				}, nil
			}),
		}

		r := httptest.NewRequest(http.MethodGet, "/og?title=test", nil)
		w := httptest.NewRecorder()

		Controller.getOg(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		if w.Header().Get("Content-Type") != "image/png" {
			t.Error("expected image/png content type")
		}

		// Should NOT have cache control in non-prod
		if w.Header().Get("Cache-Control") != "" {
			t.Error("expected no cache-control header")
		}

		if w.Body.String() != "png-data" {
			t.Error("expected body copy")
		}
	})

	t.Run("Success Production (With Cache)", func(t *testing.T) {
		config.Config().IsProduction = true

		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("png-data")),
				}, nil
			}),
		}

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		Controller.getOg(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if !strings.Contains(w.Header().Get("Cache-Control"), "max-age=31536000") {
			t.Errorf("expected aggressive cache control, got %s", w.Header().Get("Cache-Control"))
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("fetch fail")
			}),
		}

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		Controller.getOg(w, r)

		// api.HandleHttpError usually returns 500 or mapped status
		if w.Code == http.StatusOK {
			t.Error("expected error status")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		config.Config().IsProduction = false
		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("data")),
				}, nil
			}),
		}

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseWriter{ResponseWriter: w}

		Controller.getOg(errWriter, r)

		// Should log warning and exit safely
	})

	t.Run("Stream Close Error", func(t *testing.T) {
		Service.httpClient = &http.Client{
			Transport: test.CustomTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("data")},
				}, nil
			}),
		}

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		Controller.getOg(w, r)

		// Should log error in defer, but request succeeds
		if w.Code != http.StatusOK {
			t.Error("expected success despite close error")
		}
	})
}
