package og

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestController_getOg(t *testing.T) {
	t.Run("Success Default (No Cache)", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("title") != "test" {
				t.Errorf("got %s, want query param title=test", r.URL.Query().Encode())
			}
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("png-data"))

			if err != nil {
				t.Fatalf("failed to write response body: %v", err)
			}
		}))

		defer mockServer.Close()

		svc := NewService(ServiceConfig{
			OgUrl:      mockServer.URL,
			HttpClient: mockServer.Client(),
		})

		// Inject false for isProduction
		ctrl := NewController(svc, Config{ControllerConfig{IsProduction: false}})

		r := httptest.NewRequest(http.MethodGet, "/og?title=test", nil)
		w := httptest.NewRecorder()

		ctrl.getOg(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}

		if w.Header().Get("Content-Type") != "image/png" {
			t.Error("want image/png content type")
		}

		if w.Header().Get("Cache-Control") != "" {
			t.Error("want no cache-control header in non-prod")
		}

		if w.Body.String() != "png-data" {
			t.Error("want body copy")
		}
	})

	t.Run("Success Production (With Cache)", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("png-data"))

			if err != nil {
				t.Fatalf("failed to write response body: %v", err)
			}
		}))

		defer mockServer.Close()

		svc := NewService(ServiceConfig{
			OgUrl:      mockServer.URL,
			HttpClient: mockServer.Client(),
		})

		// Inject true for isProduction
		ctrl := NewController(svc, Config{ControllerConfig: ControllerConfig{IsProduction: true}})

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		ctrl.getOg(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}

		if !strings.Contains(w.Header().Get("Cache-Control"), "max-age=31536000") {
			t.Errorf("got %s, want aggressive cache control", w.Header().Get("Cache-Control"))
		}
	})

	t.Run("Service Error (Upstream 500)", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		svc := NewService(ServiceConfig{
			OgUrl:      mockServer.URL,
			HttpClient: mockServer.Client(),
		})

		ctrl := NewController(svc, Config{ControllerConfig: ControllerConfig{}})

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		ctrl.getOg(w, r)

		// The controller wraps the error using api.HandleHttpError.
		// We just check it's not OK.
		if w.Code == http.StatusOK {
			t.Error("want error status for upstream failure")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("data"))
		}))

		defer mockServer.Close()

		svc := NewService(ServiceConfig{
			OgUrl:      mockServer.URL,
			HttpClient: mockServer.Client(),
		})

		ctrl := NewController(svc, Config{ControllerConfig{}})

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseWriter{ResponseWriter: w}
		ctrl.getOg(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("Stream Close Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{
			OgUrl: "http://example.com",
			HttpClient: &http.Client{
				Transport: test.CustomTransport(func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &test.ErrorBodyCloser{Reader: strings.NewReader("data")},
					}, nil
				}),
			},
		})

		ctrl := NewController(svc, Config{ControllerConfig{}})

		r := httptest.NewRequest(http.MethodGet, "/og", nil)
		w := httptest.NewRecorder()

		ctrl.getOg(w, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Error("want success despite close error")
		}
	})
}
