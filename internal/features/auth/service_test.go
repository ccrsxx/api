package auth_test

import (
	"context"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
)

func TestService_getAuthorizationFromBearerToken(t *testing.T) {
	ctx := context.Background()

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
			svc := auth.NewService(auth.ServiceConfig{SecretKey: "test-secret"})

			_, err := svc.GetAuthorizationFromBearerToken(ctx, tt.headerToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthorizationFromBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_getAuthorizationFromQuery(t *testing.T) {
	ctx := context.Background()

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
			svc := auth.NewService(auth.ServiceConfig{SecretKey: "test-secret"})

			_, err := svc.GetAuthorizationFromQuery(ctx, tt.queryToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthorizationFromQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_getAuthorizationFromBearerOrQuery(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		headerToken string
		queryToken  string
		wantErr     bool
	}{
		{
			name:        "Valid Bearer Token (Header Preferred)",
			headerToken: "Bearer test-secret",
			queryToken:  "",
			wantErr:     false,
		},
		{
			name:        "Valid Query Token (Fallback)",
			headerToken: "",
			queryToken:  "test-secret",
			wantErr:     false,
		},
		{
			name:        "Valid Bearer Takes Priority Over Query",
			headerToken: "Bearer test-secret",
			queryToken:  "wrong-secret",
			wantErr:     false,
		},
		{
			name:        "Invalid Bearer Token (No Fallback to Query)",
			headerToken: "Bearer wrong-secret",
			queryToken:  "test-secret",
			wantErr:     true,
		},
		{
			name:        "Both Empty",
			headerToken: "",
			queryToken:  "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := auth.NewService(auth.ServiceConfig{SecretKey: "test-secret"})

			_, err := svc.GetAuthorizationFromBearerOrQuery(ctx, tt.headerToken, tt.queryToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthorizationFromBearerOrQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
