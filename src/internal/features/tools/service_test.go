package tools

import (
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/src/internal/api"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

func TestService_getIpInfo(t *testing.T) {
	originalFetcher := Service.fetcher

	defer func() {
		Service.fetcher = originalFetcher
	}()

	mockFetcher := func(ip net.IP) (*ipinfoLib.Core, error) {
		if ip.String() == "8.8.8.8" {
			return &ipinfoLib.Core{IP: net.ParseIP("8.8.8.8"), City: "Mountain View"}, nil
		}

		if ip.String() == "1.1.1.1" {
			return nil, errors.New("mock network error")
		}

		return nil, errors.New("unknown ip")
	}

	Service.fetcher = mockFetcher

	tests := []struct {
		name       string
		queryIp    string
		requestIp  string
		wantError  bool
		wantStatus int
	}{
		{
			name:      "Success with Query IP",
			queryIp:   "8.8.8.8",
			requestIp: "127.0.0.1",
			wantError: false,
		},
		{
			name:      "Success with Request IP",
			queryIp:   "",
			requestIp: "8.8.8.8",
			wantError: false,
		},
		{
			name:       "Invalid Query IP",
			queryIp:    "invalid-ip",
			requestIp:  "8.8.8.8",
			wantError:  true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid Request IP (fallback)",
			queryIp:    "",
			requestIp:  "invalid-ip",
			wantError:  true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fetcher Error",
			queryIp:    "1.1.1.1",
			requestIp:  "127.0.0.1",
			wantError:  true,
			wantStatus: 0, // General error, not HttpError
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := Service.getIpInfo(tt.queryIp, tt.requestIp)

			if tt.wantError {
				if err == nil {
					t.Error("got nil, want error")
					return
				}

				if tt.wantStatus != 0 {
					var httpErr *api.HttpError

					if errors.As(err, &httpErr) {
						if httpErr.StatusCode != tt.wantStatus {
							t.Errorf("got status %d, want %d", httpErr.StatusCode, tt.wantStatus)
						}
					} else {
						t.Error("want HttpError type")
					}
				}

				return
			}

			// Success case

			if err != nil {
				t.Errorf("unwanted error: %v", err)
			}

			if info == nil {
				t.Error("got nil, want info")
			}
		})
	}
}
