package tools

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/test"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

func TestController_GetIPAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := NewService(ServiceConfig{Fetcher: nil})
		ctrl := NewController(svc)

		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/ip", nil)

		r.RemoteAddr = "192.0.2.1:1234"

		ctrl.GetIPAddress(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}

		if w.Body.String() != "192.0.2.1" {
			t.Errorf("got %q, want 192.0.2.1", w.Body.String())
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{Fetcher: nil})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/ip", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetIPAddress(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_GetIPInfo(t *testing.T) {
	mockFetcher := func(ip net.IP) (*ipinfoLib.Core, error) {
		if ip.String() == "8.8.8.8" {
			return &ipinfoLib.Core{IP: net.ParseIP("8.8.8.8")}, nil
		}

		return nil, errors.New("mock error")
	}

	t.Run("Success", func(t *testing.T) {
		ctrl := NewController(NewService(ServiceConfig{Fetcher: mockFetcher}))

		r := httptest.NewRequest(http.MethodGet, "/ipinfo?ip=8.8.8.8", nil)
		w := httptest.NewRecorder()

		ctrl.GetIPInfo(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}

		var res map[string]any

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if _, ok := res["data"]; !ok {
			t.Fatal("want struct to be wrapped in 'data' field")
		}

		data := res["data"].(map[string]any)

		if data["ip"] != "8.8.8.8" {
			t.Errorf("got %v, want ip 8.8.8.8", data["ip"])
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{Fetcher: mockFetcher})
		ctrl := NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ipinfo?ip=1.1.1.1", nil)

		ctrl.GetIPInfo(w, r)

		if w.Code == http.StatusOK {
			t.Error("got 200, want error status")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{Fetcher: mockFetcher})
		ctrl := NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ipinfo?ip=8.8.8.8", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetIPInfo(errWriter, r)

		// The handler should have attempted to write a 200 before the write failed.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_GetHTTPHeaders(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := NewService(ServiceConfig{Fetcher: nil})
		ctrl := NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/headers", nil)
		r.Header.Set("User-Agent", "Test-Agent")

		ctrl.GetHTTPHeaders(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}

		var res map[string]string

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if val := res["User-Agent"]; val != "Test-Agent" {
			t.Errorf("got %q, want User-Agent 'Test-Agent'", val)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		svc := NewService(ServiceConfig{Fetcher: nil})
		ctrl := NewController(svc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/headers", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetHTTPHeaders(errWriter, r)

		// Confirm the handler attempted to write OK prior to the forced write error.
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
