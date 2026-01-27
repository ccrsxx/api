package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/ccrsxx/api/src/internal/api"
	"github.com/ccrsxx/api/src/internal/config"
)

type service struct{}

var Service = &service{}

func (s *service) getAuthorizationFromBearerToken(ctx context.Context, headerToken string) (string, error) {
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

	return parts[1], nil
}

func (s *service) getAuthorizationFromQuery(ctx context.Context, queryToken string) (string, error) {
	if queryToken == "" {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	match := queryToken == config.Env().SecretKey

	if !match {
		return "", &api.HttpError{
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return queryToken, nil
}
