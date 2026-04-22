package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/github"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

type beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type querier interface {
	GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error)
	GetUserWithAccountByID(ctx context.Context, id pgtype.UUID) (sqlc.GetUserWithAccountByIDRow, error)
	GetAccountByProvider(ctx context.Context, arg sqlc.GetAccountByProviderParams) (sqlc.Account, error)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	CreateAccount(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.Account, error)
	WithTx(tx pgx.Tx) querier
}

type AuthDatabaseWrapper struct {
	*sqlc.Queries
}

func (w *AuthDatabaseWrapper) WithTx(tx pgx.Tx) querier {
	return &AuthDatabaseWrapper{w.Queries.WithTx(tx)}
}

type githubClient interface {
	GetCurrentUser(ctx context.Context, accessToken string) (github.User, error)
}

type Service struct {
	db                querier
	pool              beginner
	secretKey         string
	jwtSecret         string
	githubClient      githubClient
	frontendBaseURL   string
	frontendPublicURL string
	githubOauthConfig *oauth2.Config
}

type ServiceConfig struct {
	Pool              beginner
	Database          querier
	SecretKey         string
	JwtSecret         string
	GithubClient      githubClient
	FrontendBaseURL   string
	FrontendPublicURL string
	GithubOauthConfig *oauth2.Config
}

const (
	oauthStateExpiry = 5 * time.Minute    // 5 minutes
	oauthTokenExpiry = 7 * 24 * time.Hour // 7 days / 1 week
)

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db:                cfg.Database,
		pool:              cfg.Pool,
		secretKey:         cfg.SecretKey,
		jwtSecret:         cfg.JwtSecret,
		githubClient:      cfg.GithubClient,
		frontendBaseURL:   cfg.FrontendBaseURL,
		frontendPublicURL: cfg.FrontendPublicURL,
		githubOauthConfig: cfg.GithubOauthConfig,
	}
}

func (s *Service) GetAuthorizationFromBearerToken(_ context.Context, headerToken string) (string, error) {
	if headerToken == "" {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	parts := strings.SplitN(headerToken, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token := parts[1]

	if token != s.secretKey {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return token, nil
}

func (s *Service) GetAuthorizationFromQuery(_ context.Context, queryToken string) (string, error) {
	if queryToken == "" {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	if queryToken != s.secretKey {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return queryToken, nil
}

func (s *Service) IsAdminFromOauth(ctx context.Context) (bool, error) {
	user, err := GetUserFromContext(ctx)

	if err != nil {
		return false, fmt.Errorf("is admin user from context error: %w", err)
	}

	if user.Role != "admin" {
		return false, &api.HTTPError{
			Message:    "Forbidden",
			StatusCode: http.StatusForbidden,
		}
	}

	return true, nil
}
