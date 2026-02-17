package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/src/internal/api"
	"github.com/ccrsxx/api/src/internal/config"
)

func TestService_getAuthorizationFromBearerToken(t *testing.T) {
	ctx := context.Background()

	originalKey := config.Env().SecretKey

	defer func() {
		config.Env().SecretKey = originalKey
	}()

	config.Env().SecretKey = "test-secret"

	tests := []struct {
		name        string
		headerToken string
		wantErr     bool
	}{
		{
			name:        "Valid Bearer Token",
			headerToken: "Bearer test-secret",
			wantErr:     false,
		},
		{
			name:        "Valid Lowercase Bearer",
			headerToken: "bearer test-secret",
			wantErr:     false,
		},
		{
			name:        "Invalid Secret (Wrong Token)",
			headerToken: "Bearer wrong-secret",
			wantErr:     true,
		},
		{
			name:        "Empty Header",
			headerToken: "",
			wantErr:     true,
		},
		{
			name:        "Malformed - Missing Token",
			headerToken: "Bearer",
			wantErr:     true,
		},
		{
			name:        "Malformed - Missing Bearer Prefix",
			headerToken: "test-secret",
			wantErr:     true,
		},
		{
			name:        "Malformed - Wrong Prefix",
			headerToken: "Basic test-secret",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We discard the returned string because we trust validation logic resides in the error check
			_, err := Service.getAuthorizationFromBearerToken(ctx, tt.headerToken)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getAuthorizationFromBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Optional: Ensure the error is the correct 401 type if we expect an error
			if tt.wantErr && err != nil {
				if httpErr, ok := err.(*api.HttpError); !ok || httpErr.StatusCode != http.StatusUnauthorized {
					t.Errorf("want 401 HttpError, got %v", err)
				}
			}
		})
	}
}

func TestService_getAuthorizationFromQuery(t *testing.T) {
	ctx := context.Background()

	originalKey := config.Env().SecretKey

	defer func() {
		config.Env().SecretKey = originalKey
	}()

	config.Env().SecretKey = "test-secret"

	tests := []struct {
		name       string
		queryToken string
		wantErr    bool
	}{
		{
			name:       "Valid Token",
			queryToken: "test-secret",
			wantErr:    false,
		},
		{
			name:       "Invalid Token",
			queryToken: "wrong-secret",
			wantErr:    true,
		},
		{
			name:       "Empty Token",
			queryToken: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Service.getAuthorizationFromQuery(ctx, tt.queryToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthorizationFromQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
