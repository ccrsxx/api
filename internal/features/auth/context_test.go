package auth_test

import (
	"context"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
)

func TestGetUserFromContext(t *testing.T) {
	validUser := sqlc.GetUserWithAccountByIDRow{
		Name: "test-user",
		Role: "admin",
	}

	t.Run("Valid User in Context", func(t *testing.T) {
		ctx := auth.SetUserContext(context.Background(), validUser)

		got, err := auth.GetUserFromContext(ctx)

		if err != nil {
			t.Fatalf("GetUserFromContext() unexpected error: %v", err)
		}

		if got.Name != validUser.Name {
			t.Errorf("got Name %q, want %q", got.Name, validUser.Name)
		}

		if got.Role != validUser.Role {
			t.Errorf("got Role %q, want %q", got.Role, validUser.Role)
		}
	})

	t.Run("Empty Context (No User Set)", func(t *testing.T) {
		_, err := auth.GetUserFromContext(context.Background())

		if err == nil {
			t.Fatal("GetUserFromContext() expected error, got nil")
		}
	})
}
