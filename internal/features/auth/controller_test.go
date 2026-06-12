package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/github"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/test"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

func TestController_GetCurrentUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})
		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctx := auth.SetUserContext(r.Context(), mockUserWithAccount)
		r = r.WithContext(ctx)

		w := httptest.NewRecorder()

		ctrl.GetCurrentUser(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[model.User]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Name != mockUserWithAccount.Name {
			t.Errorf("got %q, want %q", res.Data.Name, mockUserWithAccount.Name)
		}
	})

	t.Run("No Context User", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})
		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		ctrl.GetCurrentUser(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})
		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		ctx := auth.SetUserContext(r.Context(), mockUserWithAccount)
		r = r.WithContext(ctx)

		w := httptest.NewRecorder()
		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetCurrentUser(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_LoginGithub(t *testing.T) {
	t.Run("Already Logged In", func(t *testing.T) {
		jwtSecret := "controller-jwt-secret"
		userID := uuid.New()

		db := &auth.TestMockQuerier{
			MockQuerier: test.MockQuerier{
				GetUserWithAccountByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error) {
					return mockUserWithAccount, nil
				},
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret:         jwtSecret,
			FrontendPublicURL: "https://example.com",
			Database:          db,
			GithubOauthConfig: &oauth2.Config{},
		})

		token, err := svc.GenerateOauthToken(userID.String())

		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: token})

		w := httptest.NewRecorder()

		ctrl.LoginGithub(w, r)

		if w.Code != http.StatusTemporaryRedirect {
			t.Fatalf("got %d, want 307", w.Code)
		}

		location := w.Header().Get("Location")
		if location != "https://example.com/guestbook" {
			t.Errorf("got %q, want https://example.com/guestbook", location)
		}
	})

	t.Run("Not Logged In", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret: "controller-jwt-secret",
			GithubOauthConfig: &oauth2.Config{
				ClientID: "test-client-id",
				Endpoint: oauth2.Endpoint{
					AuthURL: "https://github.com/login/oauth/authorize",
				},
			},
		})

		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		ctrl.LoginGithub(w, r)

		if w.Code != http.StatusTemporaryRedirect {
			t.Fatalf("got %d, want 307", w.Code)
		}

		cookies := w.Result().Cookies()

		var found bool

		for _, c := range cookies {
			if c.Name == "oauth-state" {
				found = true

				if c.Value == "" {
					t.Error("oauth-state cookie should not be empty")
				}
			}
		}

		if !found {
			t.Error("expected oauth-state cookie to be set")
		}
	})
}

func TestController_LogoutGithub(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{
			FrontendBaseURL: "example.com",
		})

		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodPost, "/", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-token", Value: "test-token"})
		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		w := httptest.NewRecorder()

		ctrl.LogoutGithub(w, r)

		if w.Code != http.StatusNoContent {
			t.Fatalf("got %d, want 204", w.Code)
		}

		cookies := w.Result().Cookies()

		var oauthToken, oauthState *http.Cookie

		for _, c := range cookies {
			switch c.Name {
			case "oauth-token":
				oauthToken = c
			case "oauth-state":
				oauthState = c
			}
		}

		if oauthToken == nil {
			t.Fatal("expected oauth-token cookie")
		}

		if oauthToken.MaxAge != -1 {
			t.Errorf("oauth-token MaxAge got %d, want -1", oauthToken.MaxAge)
		}

		if oauthState == nil {
			t.Fatal("expected oauth-state cookie")
		}

		if oauthState.MaxAge != -1 {
			t.Errorf("oauth-state MaxAge got %d, want -1", oauthState.MaxAge)
		}
	})
}

func TestController_LoginGithubCallback(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			err := json.NewEncoder(w).Encode(map[string]any{
				"access_token": "github-access-token",
				"token_type":   "bearer",
				"expires_in":   3600,
			})

			if err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}))

		defer ts.Close()

		mockTx := newMockTxSuccess()

		db := &auth.TestMockQuerier{
			MockQuerier: test.MockQuerier{
				GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
					return sqlc.Account{}, pgx.ErrNoRows
				},
				CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
					return mockSqlcUser, nil
				},
				CreateAccountFn: func(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error) {
					return sqlc.Account{}, nil
				},
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret:         "controller-jwt-secret",
			FrontendBaseURL:   "example.com",
			FrontendPublicURL: "https://example.com",
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			GithubOauthConfig: &oauth2.Config{
				Endpoint: oauth2.Endpoint{TokenURL: ts.URL},
			},
			Database: db,
			Pool: &test.MockBeginner{
				BeginFn: func(ctx context.Context) (pgx.Tx, error) {
					return mockTx, nil
				},
			},
		})

		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state&code=test-code", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		w := httptest.NewRecorder()

		ctrl.LoginGithubCallback(w, r)

		if w.Code != http.StatusTemporaryRedirect {
			t.Fatalf("got %d, want 307", w.Code)
		}

		cookies := w.Result().Cookies()

		var found bool

		for _, c := range cookies {
			if c.Name == "oauth-token" {
				found = true

				if c.Value == "" {
					t.Error("oauth-token cookie should not be empty")
				}
			}
		}

		if !found {
			t.Error("expected oauth-token cookie to be set")
		}

		location := w.Header().Get("Location")

		if location != "https://example.com/guestbook" {
			t.Errorf("got %q, want https://example.com/guestbook", location)
		}
	})

	t.Run("ValidateOauthState Error (missing oauth state cookie)", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})
		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state&code=test-code", nil)
		w := httptest.NewRecorder()

		ctrl.LoginGithubCallback(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("CreateOauthTokenForGithubUser Error (github client error)", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			err := json.NewEncoder(w).Encode(map[string]any{
				"access_token": "github-access-token",
				"token_type":   "bearer",
				"expires_in":   3600,
			})

			if err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}

		}))

		defer ts.Close()

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return github.User{}, errors.New("github error")
				},
			},
			GithubOauthConfig: &oauth2.Config{
				Endpoint: oauth2.Endpoint{TokenURL: ts.URL},
			},
		})

		ctrl := auth.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state&code=test-code", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		w := httptest.NewRecorder()

		ctrl.LoginGithubCallback(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})
}
