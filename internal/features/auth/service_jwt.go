package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// signToken is a package-level function that can be overridden in tests
// to simulate signing errors.
var signToken = func(token *jwt.Token, key []byte) (string, error) {
	return token.SignedString(key)
}

// parseToken is a package-level function that can be overridden in tests
// to simulate parse errors with unexpected claims types.
var parseToken = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, opts ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, keyFunc, opts...)
}

func (s *Service) GenerateOauthToken(userID string) (string, error) {
	currentTime := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    "API Portofolio",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(oauthTokenExpiry)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := signToken(token, []byte(s.jwtSecret))

	if err != nil {
		return "", fmt.Errorf("jwt sign token error: %w", err)
	}

	return ss, nil
}

func (s *Service) ValidateOauthToken(ctx context.Context, r *http.Request) (sqlc.GetUserWithAccountByIDRow, error) {
	oauthToken, err := r.Cookie("oauth-token")

	// r.Cookie returns always return nil or ErrNoCookie
	// Checking just generic error is enough to assume no cookie found
	if err != nil {
		slog.Warn("jwt validate missing token error", "error", err)

		return sqlc.GetUserWithAccountByIDRow{}, &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token, err := parseToken(
		oauthToken.Value,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) { return []byte(s.jwtSecret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		slog.Warn("jwt validate parse token error", "error", err)

		return sqlc.GetUserWithAccountByIDRow{}, &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok {
		return sqlc.GetUserWithAccountByIDRow{}, errors.New("jwt validate unknown claim error")
	}

	id, err := uuid.Parse(claims.Subject)

	if err != nil {
		return sqlc.GetUserWithAccountByIDRow{}, fmt.Errorf("jwt validate user id error: %w", err)
	}

	user, err := s.db.GetUserWithAccountByID(ctx, pgtype.UUID{Bytes: id, Valid: true})

	if err != nil {
		return sqlc.GetUserWithAccountByIDRow{}, fmt.Errorf("jwt validate get user error: %w", err)
	}

	return user, nil
}
