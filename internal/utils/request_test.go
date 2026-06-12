package utils_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/utils"
)

func TestGetIPAddressFromRequest(t *testing.T) {
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

			setHTTPHeaders(r, tt.headers)

			r.RemoteAddr = tt.remoteAddr

			got := utils.GetIPAddressFromRequest(r)

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "Loopback IPv4",
			ip:   "127.0.0.1",
			want: true,
		},
		{
			name: "Loopback IPv6",
			ip:   "::1",
			want: true,
		},
		{
			name: "Private 10.x.x.x",
			ip:   "10.0.0.1",
			want: true,
		},
		{
			name: "Private 172.16.x.x",
			ip:   "172.16.0.1",
			want: true,
		},
		{
			name: "Private 192.168.x.x",
			ip:   "192.168.1.100",
			want: true,
		},
		{
			name: "Public IP",
			ip:   "8.8.8.8",
			want: false,
		},
		{
			name: "Public 172.32.x.x (Outside Private Range)",
			ip:   "172.32.0.1",
			want: false,
		},
		{
			name: "Invalid IP",
			ip:   "not-an-ip",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsPrivateIP(tt.ip)

			if got != tt.want {
				t.Errorf("IsPrivateIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestGetHttpHeadersFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    map[string]string
	}{
		{
			name:    "Single Headers",
			headers: map[string]string{"Content-Type": "application/json"},
			want:    map[string]string{"Content-Type": "application/json"},
		},
		{
			name:    "Multiple Values for Same Header",
			headers: map[string]string{"X-Multi": "Value1, Value2"},
			want:    map[string]string{"X-Multi": "Value1, Value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			setHTTPHeaders(r, tt.headers)

			got := utils.GetHTTPHeadersFromRequest(r)

			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("got %v, want %v", got[k], v)
				}
			}
		})
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

			setHTTPHeaders(r, tt.headers)

			if tt.isTLS {
				r.TLS = &tls.ConnectionState{}
			}

			got := utils.GetPublicURLFromRequest(r)

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func setHTTPHeaders(r *http.Request, headers map[string]string) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}
