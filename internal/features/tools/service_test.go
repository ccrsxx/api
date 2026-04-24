package tools_test

import (
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/features/tools"
	ipinfoLib "github.com/ipinfo/go/v2/ipinfo"
)

func TestService_GetIPInfo(t *testing.T) {
	mock := &mockIPInfoClient{
		result: func(ip net.IP) (*ipinfoLib.Core, error) {
			if ip.String() == "8.8.8.8" {
				return &ipinfoLib.Core{IP: net.ParseIP("8.8.8.8"), City: "Mountain View"}, nil
			}

			if ip.String() == "1.1.1.1" {
				return nil, errors.New("mock network error")
			}

			return nil, errors.New("unknown ip")
		},
	}

	tests := []struct {
		name       string
		queryIP    string
		requestIP  string
		wantError  bool
		wantStatus int
	}{
		{
			name:      "Success with Query IP",
			queryIP:   "8.8.8.8",
			requestIP: "127.0.0.1",
			wantError: false,
		},
		{
			name:      "Success with Request IP",
			queryIP:   "",
			requestIP: "8.8.8.8",
			wantError: false,
		},
		{
			name:       "Invalid Query IP",
			queryIP:    "invalid-ip",
			requestIP:  "8.8.8.8",
			wantError:  true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid Request IP (fallback)",
			queryIP:    "",
			requestIP:  "invalid-ip",
			wantError:  true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fetcher Error",
			queryIP:    "1.1.1.1",
			requestIP:  "127.0.0.1",
			wantError:  true,
			wantStatus: 0, // General error, not HttpError
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tools.NewService(tools.ServiceConfig{IPInfoClient: mock})

			info, err := svc.GetIPInfo(tt.queryIP, tt.requestIP)

			if tt.wantError {
				if err == nil {
					t.Error("got nil, want error")
					return
				}

				if tt.wantStatus != 0 {
					if httpErr, ok := errors.AsType[*api.HTTPError](err); ok {
						if httpErr.StatusCode != tt.wantStatus {
							t.Errorf("got status %d, want %d", httpErr.StatusCode, tt.wantStatus)
						}
					} else {
						t.Error("want HTTPError type")
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
