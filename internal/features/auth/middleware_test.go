package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestMiddleware_IsAuthorizedFromBearer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "middleware-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer middleware-secret")

		mw.IsAuthorizedFromBearer(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "middleware-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer wrong-key")

		mw.IsAuthorizedFromBearer(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}

func TestMiddleware_IsAuthorizedFromQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "query-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=query-secret", nil)

		mw.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "query-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)

		mw.IsAuthorizedFromQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}

func TestMiddleware_IsAuthorizedFromBearerOrQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success via Bearer Header (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "combo-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer combo-secret")

		mw.IsAuthorizedFromBearerOrQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Success via Query Token (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "combo-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=combo-secret", nil)

		mw.IsAuthorizedFromBearerOrQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want status 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{SecretKey: "combo-secret"}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?token=wrong-secret", nil)
		r.Header.Set("Authorization", "Bearer wrong-secret")

		mw.IsAuthorizedFromBearerOrQuery(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want status 401", w.Code)
		}
	})
}

func TestMiddleware_IsAuthorizedFromOauth(t *testing.T) {
	t.Run("Success (200 OK)", func(t *testing.T) {
		jwtSecret := "middleware-jwt-secret"
		userID := uuid.New()

		db := &auth.TestMockQuerier{
			MockQuerier: test.MockQuerier{
				GetUserWithAccountByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error) {
					return sqlc.GetUserWithAccountByIDRow{
						Name: "Test User",
						Role: "admin",
					}, nil
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

		mw := auth.NewMiddleware(svc)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.GetUserFromContext(r.Context())

			if err != nil {
				t.Errorf("user not in context: %v", err)
			}

			if user.Name != "Test User" {
				t.Errorf("got %q, want Test User", user.Name)
			}

			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: token})

		mw.IsAuthorizedFromOauth(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}
	})

	t.Run("Fail (401 Unauthorized)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{JwtSecret: "middleware-jwt-secret"}))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		mw.IsAuthorizedFromOauth(handler).ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", w.Code)
		}
	})
}

func TestMiddleware_IsAdminFromOauth(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Success (200 OK)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctx := auth.SetUserContext(r.Context(), sqlc.GetUserWithAccountByIDRow{
			Name: "Admin User",
			Role: "admin",
		})

		r = r.WithContext(ctx)

		mw.IsAdminFromOauth(handler).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want 200", w.Code)
		}
	})

	t.Run("Fail (403 Forbidden)", func(t *testing.T) {
		mw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctx := auth.SetUserContext(r.Context(), sqlc.GetUserWithAccountByIDRow{
			Name: "Regular User",
			Role: "user",
		})

		r = r.WithContext(ctx)

		mw.IsAdminFromOauth(handler).ServeHTTP(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("got %d, want 403", w.Code)
		}
	})
}
