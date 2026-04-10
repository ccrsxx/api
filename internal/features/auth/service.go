package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/ccrsxx/api/internal/api"
)

type Service struct {
	secretKey string
}

type ServiceConfig struct {
	SecretKey string
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		secretKey: cfg.SecretKey,
	}
}

func (s *Service) getAuthorizationFromBearerToken(ctx context.Context, headerToken string) (string, error) {
	if headerToken == "" {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	parts := strings.SplitN(headerToken, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	token := parts[1]

	if token != s.secretKey {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return token, nil
}

func (s *Service) getAuthorizationFromQuery(ctx context.Context, queryToken string) (string, error) {
	if queryToken == "" {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	if queryToken != s.secretKey {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return queryToken, nil
}
