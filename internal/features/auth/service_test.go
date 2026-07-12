package auth_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
)

func TestNewService_BindsRealQuerierToTx(t *testing.T) {
	base := &sqlc.Queries{}

	svc := auth.NewService(auth.ServiceConfig{Database: base})

	bound := svc.NewTxQuerier(&test.MockTx{})

	// A real *sqlc.Queries must be bound to the tx (a fresh instance), not
	// reused as-is (which would run queries on the pool, outside the tx).
	if bound == base {
		t.Fatal("safety net did not bind to the tx: got the pool-bound base instance")
	}
}

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

func TestService_IsAdminFromOauth(t *testing.T) {
	t.Run("Success (Admin)", func(t *testing.T) {
		ctx := auth.SetUserContext(context.Background(), sqlc.GetUserWithAccountByIDRow{
			Name: "Admin User",
			Role: "admin",
		})

		svc := auth.NewService(auth.ServiceConfig{})

		isAdmin, err := svc.IsAdminFromOauth(ctx)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if !isAdmin {
			t.Error("expected true, got false")
		}
	})

	t.Run("Non-Admin (403 Forbidden)", func(t *testing.T) {
		ctx := auth.SetUserContext(context.Background(), sqlc.GetUserWithAccountByIDRow{
			Name: "Regular User",
			Role: "user",
		})

		svc := auth.NewService(auth.ServiceConfig{})

		_, err := svc.IsAdminFromOauth(ctx)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want *api.HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusForbidden {
			t.Errorf("got %d, want 403", httpErr.StatusCode)
		}
	})

	t.Run("No User in Context", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})

		_, err := svc.IsAdminFromOauth(context.Background())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "is admin user from context error") {
			t.Errorf("got %v, want is admin user from context error", err)
		}
	})
}
