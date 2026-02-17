package tools

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/src/internal/test"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

func TestController_GetIpAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/ip", nil)
		r.RemoteAddr = "192.0.2.1:1234"

		Controller.GetIpAddress(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		if w.Body.String() != "192.0.2.1" {
			t.Errorf("want 192.0.2.1, got %q", w.Body.String())
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ip", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		Controller.GetIpAddress(errWriter, r)
	})
}

func TestController_GetIpInfo(t *testing.T) {
	originalFetcher := Service.fetcher

	defer func() { Service.fetcher = originalFetcher }()

	Service.fetcher = func(ip net.IP) (*ipinfoLib.Core, error) {
		if ip.String() == "8.8.8.8" {
			return &ipinfoLib.Core{IP: net.ParseIP("8.8.8.8")}, nil
		}

		return nil, errors.New("mock error")
	}

	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ipinfo?ip=8.8.8.8", nil)
		w := httptest.NewRecorder()

		Controller.GetIpInfo(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
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
			t.Errorf("want ip 8.8.8.8, got %v", data["ip"])
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ipinfo?ip=1.1.1.1", nil)

		Controller.GetIpInfo(w, r)

		if w.Code == http.StatusOK {
			t.Error("want error status, got 200")
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ipinfo?ip=8.8.8.8", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		Controller.GetIpInfo(errWriter, r)
	})
}

func TestController_GetHttpHeaders(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/headers", nil)

		r.Header.Set("User-Agent", "Test-Agent")

		Controller.GetHttpHeaders(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want 200, got %d", w.Code)
		}

		var res map[string]string

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if val := res["User-Agent"]; val != "Test-Agent" {
			t.Errorf("want User-Agent 'Test-Agent', got %q", val)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/headers", nil)

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		Controller.GetHttpHeaders(errWriter, r)
	})
}
