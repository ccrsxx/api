package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/config"
)

type service struct{}

var Service = &service{}

func (s *service) getAuthorizationFromBearerToken(_ context.Context, headerToken string) (string, error) {
	if headerToken == "" {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	parts := strings.SplitN(headerToken, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	return parts[1], nil
}

func (s *service) getAuthorizationFromQuery(_ context.Context, queryToken string) (string, error) {
	if queryToken == "" {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	match := queryToken == config.Env().SecretKey

	if !match {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	return queryToken, nil
}
