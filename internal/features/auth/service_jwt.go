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

func (s *Service) GenerateOauthToken(userID string) (string, error) {
	currentTime := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    "API Portofolio",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(oauthTokenExpiry)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(s.jwtSecret))

	if err != nil {
		return "", fmt.Errorf("jwt sign token error: %w", err)
	}

	return ss, nil
}

func (s *Service) ValidateOauthToken(ctx context.Context, r *http.Request) (sqlc.GetUserWithAccountByIDRow, error) {
	oauthToken, err := r.Cookie("oauth-token")

	if err != nil {
		slog.Warn("jwt validate cookie token error", "error", err)

		return sqlc.GetUserWithAccountByIDRow{}, &api.HTTPError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token, err := jwt.ParseWithClaims(
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
