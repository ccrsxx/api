package utils

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIpAddressFromRequest(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name:       "Cloudflare Header (Priority 1)",
			headers:    map[string]string{"CF-Connecting-IP": "1.1.1.1", "X-Real-IP": "2.2.2.2"},
			remoteAddr: "127.0.0.1:1234",
			want:       "1.1.1.1",
		},
		{
			name:       "X-Real-IP Header (Priority 2)",
			headers:    map[string]string{"X-Real-IP": "2.2.2.2"},
			remoteAddr: "127.0.0.1:1234",
			want:       "2.2.2.2",
		},
		{
			name:       "X-Forwarded-For Single",
			headers:    map[string]string{"X-Forwarded-For": "3.3.3.3"},
			remoteAddr: "127.0.0.1:1234",
			want:       "3.3.3.3",
		},
		{
			name:       "X-Forwarded-For Multiple (Take First)",
			headers:    map[string]string{"X-Forwarded-For": "4.4.4.4, 5.5.5.5"},
			remoteAddr: "127.0.0.1:1234",
			want:       "4.4.4.4",
		},
		{
			name:       "No Headers (Fallback to RemoteAddr)",
			headers:    nil,
			remoteAddr: "192.168.1.50:4000",
			want:       "192.168.1.50",
		},
		{
			name:       "RemoteAddr Invalid Format (Return As Is)",
			headers:    nil,
			remoteAddr: "invalid-ip-string", // net.SplitHostPort will fail
			want:       "invalid-ip-string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			setHttpHeaders(r, tt.headers)

			r.RemoteAddr = tt.remoteAddr

			got := GetIpAddressFromRequest(r)

			if got != tt.want {
				t.Errorf("GetIpAddressFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHttpHeadersFromRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("X-Custom", "Value1")
	r.Header.Add("X-Custom", "Value2")

	got := GetHttpHeadersFromRequest(r)

	if got["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", got["Content-Type"])
	}

	if got["X-Custom"] != "Value1, Value2" {
		t.Errorf("Expected X-Custom 'Value1, Value2', got %s", got["X-Custom"])
	}
}

func TestGetPublicUrlFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		isTLS   bool
		host    string
		want    string
	}{
		{
			name:    "Standard HTTP",
			headers: nil,
			isTLS:   false,
			host:    "example.com",
			want:    "http://example.com",
		},
		{
			name:    "HTTPS via Proxy Header",
			headers: map[string]string{"X-Forwarded-Proto": "https"},
			isTLS:   false,
			host:    "api.example.com",
			want:    "https://api.example.com",
		},
		{
			name:    "Native HTTPS (TLS)",
			headers: nil,
			isTLS:   true,
			host:    "secure.com",
			want:    "https://secure.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://"+tt.host, nil)

			r.Host = tt.host

			setHttpHeaders(r, tt.headers)

			if tt.isTLS {
				r.TLS = &tls.ConnectionState{}
			}

			got := GetPublicUrlFromRequest(r)

			if got != tt.want {
				t.Errorf("GetPublicUrlFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setHttpHeaders(r *http.Request, headers map[string]string) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}
