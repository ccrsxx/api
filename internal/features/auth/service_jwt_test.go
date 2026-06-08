package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestService_GenerateOauthToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{JwtSecret: "test-jwt-secret"})

		token, err := svc.GenerateOauthToken(uuid.New().String())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if token == "" {
			t.Fatal("expected non-empty token")
		}
	})

	t.Run("Sign Error", func(t *testing.T) {
		restore := auth.SetSignToken(func(token *jwt.Token, key []byte) (string, error) {
			return "", errors.New("signing error")
		})

		defer restore()

		svc := auth.NewService(auth.ServiceConfig{JwtSecret: "test-jwt-secret"})

		_, err := svc.GenerateOauthToken(uuid.New().String())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "jwt sign token error") {
			t.Errorf("got %v, want jwt sign token error", err)
		}
	})
}

func TestService_ValidateOauthToken(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()

		db := &auth.TestMockQuerier{
			MockQuerier: test.MockQuerier{
				GetUserWithAccountByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error) {
					return mockUserWithAccount, nil
				},
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret: jwtSecret,
			Database:  db,
		})

		token, err := svc.GenerateOauthToken(userID.String())

		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: token})

		user, err := svc.ValidateOauthToken(r.Context(), r)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if user.Name != mockUserWithAccount.Name {
			t.Errorf("got %q, want %q", user.Name, mockUserWithAccount.Name)
		}
	})

	t.Run("Missing Cookie", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{JwtSecret: jwtSecret})

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		_, err := svc.ValidateOauthToken(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "Invalid token") {
			t.Errorf("got %v, want Invalid token", err)
		}
	})

	t.Run("Invalid JWT", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{JwtSecret: jwtSecret})

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: "invalid-jwt-token"})

		_, err := svc.ValidateOauthToken(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "Invalid token") {
			t.Errorf("got %v, want Invalid token", err)
		}
	})

	t.Run("Invalid UUID in Subject", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{JwtSecret: jwtSecret})

		token, err := svc.GenerateOauthToken("not-a-uuid")

		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: token})

		_, err = svc.ValidateOauthToken(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "jwt validate user id error") {
			t.Errorf("got %v, want jwt validate user id error", err)
		}
	})

	t.Run("GetUserWithAccountByID Error", func(t *testing.T) {
		userID := uuid.New()

		db := &auth.TestMockQuerier{
			MockQuerier: test.MockQuerier{
				GetUserWithAccountByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error) {
					return sqlc.GetUserWithAccountByIDRow{}, errors.New("db error")
				},
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret: jwtSecret,
			Database:  db,
		})

		token, err := svc.GenerateOauthToken(userID.String())
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: token})

		_, err = svc.ValidateOauthToken(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "jwt validate get user error") {
			t.Errorf("got %v, want jwt validate get user error", err)
		}
	})

	t.Run("Unknown Claims Type", func(t *testing.T) {
		restore := auth.SetParseToken(func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, opts ...jwt.ParserOption) (*jwt.Token, error) {
			return &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{},
			}, nil
		})

		defer restore()

		svc := auth.NewService(auth.ServiceConfig{JwtSecret: jwtSecret})

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: "any-token"})

		_, err := svc.ValidateOauthToken(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "jwt validate unknown claim error") {
			t.Errorf("got %v, want jwt validate unknown claim error", err)
		}
	})
}
