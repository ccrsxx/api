package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/clients/github"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

type mockGithubClient struct {
	GetCurrentUserFn func(ctx context.Context, accessToken string) (github.User, error)
}

func (m *mockGithubClient) GetCurrentUser(ctx context.Context, accessToken string) (github.User, error) {
	return m.GetCurrentUserFn(ctx, accessToken)
}

var (
	mockUserID   = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockSqlcUser = sqlc.User{
		ID:    mockUserID,
		Name:  "Test User",
		Role:  "admin",
		Image: pgtype.Text{String: "https://example.com/avatar.jpg", Valid: true},
	}
	mockUserWithAccount = sqlc.GetUserWithAccountByIDRow{
		ID:       mockUserID,
		Name:     "Test User",
		Role:     "admin",
		Image:    pgtype.Text{String: "https://example.com/avatar.jpg", Valid: true},
		Username: pgtype.Text{String: "testuser", Valid: true},
	}
	mockGithubUserName  = "Test User"
	mockGithubUserEmail = "test@example.com"
	mockGithubUser      = github.User{
		ID:        12345,
		Login:     "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		Name:      &mockGithubUserName,
		Email:     &mockGithubUserEmail,
	}
)

func newMockTxSuccess() *test.MockTx {
	return &test.MockTx{
		CommitFn:   func(ctx context.Context) error { return nil },
		RollbackFn: func(ctx context.Context) error { return pgx.ErrTxClosed },
	}
}

func TestService_CreateOauthTokenForGithubUser(t *testing.T) {
	t.Run("Success (New User)", func(t *testing.T) {
		mockTx := newMockTxSuccess()

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{}, pgx.ErrNoRows
			},
			CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
				return mockSqlcUser, nil
			},
			CreateAccountFn: func(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error) {
				return sqlc.Account{}, nil
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret: "test-jwt-secret",
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
			Pool: &test.MockBeginner{
				BeginFn: func(ctx context.Context) (pgx.Tx, error) {
					return mockTx, nil
				},
			},
		})

		token, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if token == "" {
			t.Fatal("expected non-empty token")
		}
	})

	t.Run("Success (Existing User, Needs Update)", func(t *testing.T) {
		newName := "Updated Name"

		githubUser := github.User{
			ID:        12345,
			Login:     "testuser",
			AvatarURL: "https://example.com/new-avatar.jpg",
			Name:      &newName,
		}

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{UserID: mockUserID}, nil
			},
			GetUserByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
				return mockSqlcUser, nil
			},
			UpdateUserFn: func(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
				return sqlc.User{ID: mockUserID, Name: newName}, nil
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret: "test-jwt-secret",
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return githubUser, nil
				},
			},
			Database: db,
		})

		token, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if token == "" {
			t.Fatal("expected non-empty token")
		}
	})

	t.Run("GetCurrentUser Error", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return github.User{}, errors.New("github api error")
				},
			},
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create session get github user error") {
			t.Errorf("got %v, want create session get github user error", err)
		}
	})

	t.Run("GetAccountByProvider DB Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{}, errors.New("db error")
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create github user get account error") {
			t.Errorf("got %v, want create github user get account error", err)
		}
	})

	t.Run("GetUserByID Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{UserID: mockUserID}, nil
			},
			GetUserByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
				return sqlc.User{}, errors.New("db error")
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "update github user get user error") {
			t.Errorf("got %v, want update github user get user error", err)
		}
	})

	t.Run("UpdateUser Error", func(t *testing.T) {
		newName := "Different Name"
		githubUser := github.User{
			ID:        12345,
			Login:     "testuser",
			AvatarURL: "https://example.com/new-avatar.jpg",
			Name:      &newName,
		}

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{UserID: mockUserID}, nil
			},
			GetUserByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
				return mockSqlcUser, nil
			},
			UpdateUserFn: func(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
				return sqlc.User{}, errors.New("db error")
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return githubUser, nil
				},
			},
			Database: db,
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "update github user update db error") {
			t.Errorf("got %v, want update github user update db error", err)
		}
	})

	t.Run("Begin Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{}, pgx.ErrNoRows
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
			Pool: &test.MockBeginner{
				BeginFn: func(ctx context.Context) (pgx.Tx, error) {
					return nil, errors.New("pool error")
				},
			},
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create new github user begin tx error") {
			t.Errorf("got %v, want create new github user begin tx error", err)
		}
	})

	t.Run("CreateAccount Error", func(t *testing.T) {
		mockTx := newMockTxSuccess()

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{}, pgx.ErrNoRows
			},
			CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
				return mockSqlcUser, nil
			},
			CreateAccountFn: func(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error) {
				return sqlc.Account{}, errors.New("db error")
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
			Pool: &test.MockBeginner{
				BeginFn: func(ctx context.Context) (pgx.Tx, error) {
					return mockTx, nil
				},
			},
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create new github account insert error") {
			t.Errorf("got %v, want create new github account insert error", err)
		}
	})

	t.Run("Commit Error", func(t *testing.T) {
		mockTx := &test.MockTx{
			CommitFn:   func(ctx context.Context) error { return errors.New("commit error") },
			RollbackFn: func(ctx context.Context) error { return pgx.ErrTxClosed },
		}

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{}, pgx.ErrNoRows
			},
			CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
				return mockSqlcUser, nil
			},
			CreateAccountFn: func(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error) {
				return sqlc.Account{}, nil
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
			Pool: &test.MockBeginner{
				BeginFn: func(ctx context.Context) (pgx.Tx, error) {
					return mockTx, nil
				},
			},
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create new github user commit tx error") {
			t.Errorf("got %v, want create new github user commit tx error", err)
		}
	})

	t.Run("Rollback Warning", func(t *testing.T) {
		mockTx := &test.MockTx{
			RollbackFn: func(ctx context.Context) error { return errors.New("rollback error") },
		}

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{}, pgx.ErrNoRows
			},
			CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
				return sqlc.User{}, errors.New("db error")
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
			Pool: &test.MockBeginner{
				BeginFn: func(ctx context.Context) (pgx.Tx, error) {
					return mockTx, nil
				},
			},
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create new github user insert error") {
			t.Errorf("got %v, want create new github user insert error", err)
		}
	})

	t.Run("GenerateOauthToken Error", func(t *testing.T) {
		restore := auth.SetSignToken(func(token *jwt.Token, key []byte) (string, error) {
			return "", errors.New("signing error")
		})

		defer restore()

		db := &test.MockQuerier{
			GetAccountByProviderFn: func(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error) {
				return sqlc.Account{UserID: mockUserID}, nil
			},
			GetUserByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
				return mockSqlcUser, nil
			},
		}

		svc := auth.NewService(auth.ServiceConfig{
			JwtSecret: "test-jwt-secret",
			GithubClient: &mockGithubClient{
				GetCurrentUserFn: func(ctx context.Context, accessToken string) (github.User, error) {
					return mockGithubUser, nil
				},
			},
			Database: db,
		})

		_, err := svc.CreateOauthTokenForGithubUser(context.Background(), "test-access-token")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create session generate session token error") {
			t.Errorf("got %v, want create session generate session token error", err)
		}
	})
}

func TestService_ValidateOauthState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			err := json.NewEncoder(w).Encode(map[string]any{
				"access_token": "test-access-token",
				"token_type":   "bearer",
				"expires_in":   3600,
			})

			if err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}))

		defer ts.Close()

		svc := auth.NewService(auth.ServiceConfig{
			GithubOauthConfig: &oauth2.Config{
				Endpoint: oauth2.Endpoint{TokenURL: ts.URL},
			},
		})

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state&code=test-code", nil)
		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		token, err := svc.ValidateOauthState(r.Context(), r)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if token.AccessToken != "test-access-token" {
			t.Errorf("got %q, want test-access-token", token.AccessToken)
		}
	})

	t.Run("Missing Cookie", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state", nil)

		_, err := svc.ValidateOauthState(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "validate oauth request cookie error") {
			t.Errorf("got %v, want validate oauth request cookie error", err)
		}
	})

	t.Run("State Mismatch", func(t *testing.T) {
		svc := auth.NewService(auth.ServiceConfig{})

		r := httptest.NewRequest(http.MethodGet, "/?state=wrong-state", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "correct-state"})

		_, err := svc.ValidateOauthState(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "validate oauth request invalid state") {
			t.Errorf("got %v, want validate oauth request invalid state", err)
		}
	})

	t.Run("Exchange Error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))

		defer ts.Close()

		svc := auth.NewService(auth.ServiceConfig{
			GithubOauthConfig: &oauth2.Config{
				Endpoint: oauth2.Endpoint{TokenURL: ts.URL},
			},
		})

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state&code=bad-code", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		_, err := svc.ValidateOauthState(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "validate oauth request exchange error") {
			t.Errorf("got %v, want validate oauth request exchange error", err)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			err := json.NewEncoder(w).Encode(map[string]any{
				"access_token": "test-token",
				"token_type":   "bearer",
				"expires_in":   -3600,
			})

			if err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}))

		defer ts.Close()

		svc := auth.NewService(auth.ServiceConfig{
			GithubOauthConfig: &oauth2.Config{
				Endpoint: oauth2.Endpoint{TokenURL: ts.URL},
			},
		})

		r := httptest.NewRequest(http.MethodGet, "/?state=test-state&code=test-code", nil)

		r.AddCookie(&http.Cookie{Name: "oauth-state", Value: "test-state"})

		_, err := svc.ValidateOauthState(r.Context(), r)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "validate oauth request invalid token") {
			t.Errorf("got %v, want validate oauth request invalid token", err)
		}
	})
}
