package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
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
}

// NewSqlcTxFactory returns a factory that binds sqlc queries to a transaction.
// It lets the service open transactions without a per-feature database wrapper:
// *sqlc.Queries already satisfies querier, and WithTx returns a tx-scoped copy.
func NewSqlcTxFactory(db *sqlc.Queries) func(pgx.Tx) querier {
	return func(tx pgx.Tx) querier {
		return db.WithTx(tx)
	}
}

type githubClient interface {
	GetCurrentUser(ctx context.Context, accessToken string) (github.User, error)
}

type Service struct {
	db                querier
	pool              beginner
	newTx             func(pgx.Tx) querier
	secretKey         string
	jwtSecret         string
	githubClient      githubClient
	frontendBaseURL   string
	frontendPublicURL string
	githubOauthConfig *oauth2.Config
}

type ServiceConfig struct {
	Pool     beginner
	Database querier
	// WithTx binds the querier to a transaction. It is optional: when nil,
	// NewService auto-binds a real *sqlc.Queries via NewSqlcTxFactory, and
	// reuses the instance for mocks. Set it only for a custom querier type.
	WithTx            func(pgx.Tx) querier
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
	newTx := cfg.WithTx

	if newTx == nil {
		if db, ok := cfg.Database.(*sqlc.Queries); ok {
			// Safety net: a real sqlc querier MUST bind queries to the
			// transaction, even when WithTx was not wired explicitly. Otherwise
			// queries would run on the pool, outside the tx, losing atomicity.
			newTx = NewSqlcTxFactory(db)
		} else {
			// Tests/mocks: the querier ignores the tx, so the same instance is
			// reused inside the transaction.
			newTx = func(pgx.Tx) querier {
				return cfg.Database
			}
		}
	}

	return &Service{
		db:                cfg.Database,
		pool:              cfg.Pool,
		newTx:             newTx,
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

	return s.validateSecretKey(token)
}

func (s *Service) GetAuthorizationFromQuery(_ context.Context, queryToken string) (string, error) {
	if queryToken == "" {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return s.validateSecretKey(queryToken)
}

func (s *Service) GetAuthorizationFromBearerOrQuery(ctx context.Context, headerToken, queryToken string) (string, error) {
	if headerToken != "" {
		return s.GetAuthorizationFromBearerToken(ctx, headerToken)
	}

	return s.GetAuthorizationFromQuery(ctx, queryToken)
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

// Compare the token against the secret in constant time to avoid leaking
// information via timing differences. We hash both values first so the
// comparison runs over fixed-length (32-byte) inputs; subtle.ConstantTimeCompare
// would otherwise short-circuit on a length mismatch and leak the secret's length.
//
// Rate limiting (Cloudflare + API Gateway) is the primary defense against
// brute-force/timing attacks; this is defense-in-depth in case those fail open.
func (s *Service) validateSecretKey(token string) (string, error) {
	tokenHash := sha256.Sum256([]byte(token))
	secretKeyHash := sha256.Sum256([]byte(s.secretKey))

	if subtle.ConstantTimeCompare(tokenHash[:], secretKeyHash[:]) != 1 {
		return "", &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return token, nil
}
