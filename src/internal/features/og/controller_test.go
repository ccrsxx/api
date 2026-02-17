package og

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/test"
)

func TestController_getOg(t *testing.T) {
	originalOgUrl := Service.ogUrl
	originalClient := Service.httpClient

	originalDev := config.Config().IsDevelopment
	originalProd := config.Config().IsProduction

	defer func() {
		Service.ogUrl = originalOgUrl
		Service.httpClient = originalClient

		config.Config().IsProduction = originalProd
		config.Config().IsDevelopment = originalDev
	}()

	// Default to Production for most tests so we can inject the mock URL
	config.Config().IsDevelopment = false

	t.Run("Success Default (No Cache)", func(t *testing.T) {
		config.Config().IsProduction = false

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("title") != "test" {
				t.Errorf("want query param title=test, got %s", r.URL.Query().Encode())
			}
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte("png-data"))

			if err != nil {
				t.Fatalf("failed to write response body: %v", err)
			}
		}))

		defer mockServer.Close()

		Service.ogUrl = mockServer.URL
		Service.httpClient = mockServer.Client()

		r := httptest.NewRequest(http.MethodGet, "/og?title=test", nil)
		w := httptest.NewRecorder()

		Controller.getOg(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		if w.Header().Get("Content-Type") != "image/png" {
			t.Error("expected image/png content type")
		}

		if w.Header().Get("Cache-Control") != "" {
			t.Error("expected no cache-control header in non-prod")
		}

		if w.Body.String() != "png-data" {
			t.Error("expected body copy")
		}
	})

	t.Run("Success Production (With Cache)", func(t *testing.T) {
		config.Config().IsProduction = true

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte("png-data"))

			if err != nil {
				t.Fatalf("failed to write response body: %v", err)
			}
		}))

		defer mockServer.Close()

		Service.ogUrl = mockServer.URL
		Service.httpClient = mockServer.Client()

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

	t.Run("Service Error (Upstream 500)", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		Service.ogUrl = mockServer.URL
		Service.httpClient = mockServer.Client()

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		Controller.getOg(w, r)

		// The controller wraps the error using api.HandleHttpError.
		// We just check it's not OK.
		if w.Code == http.StatusOK {
			t.Error("expected error status for upstream failure")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		config.Config().IsProduction = false

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte("data"))

			if err != nil {
				t.Fatalf("failed to write response body: %v", err)
			}
		}))

		defer mockServer.Close()

		Service.ogUrl = mockServer.URL
		Service.httpClient = mockServer.Client()

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseWriter{ResponseWriter: w}

		Controller.getOg(errWriter, r)

		if w.Code != http.StatusOK {
			t.Error("expected 200 even if write error occurs")
		}
	})

	t.Run("Stream Close Error", func(t *testing.T) {
		Service.ogUrl = "http://example.com"
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

		if w.Code != http.StatusOK {
			t.Error("expected success despite close error")
		}
	})
}
