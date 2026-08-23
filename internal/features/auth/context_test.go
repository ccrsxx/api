package auth_test

import (
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
		ctx := auth.SetUserContext(t.Context(), validUser)

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
		_, err := auth.GetUserFromContext(t.Context())

		if err == nil {
			t.Fatal("GetUserFromContext() expected error, got nil")
		}
	})
}
