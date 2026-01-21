package auth

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api-go/src/internal/api"
)

type service struct{}

var Service = &service{}

func (s *service) getAuthorizationFromBearerToken(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")

	if authorization == "" {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	parts := strings.SplitN(authorization, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	return parts[1], nil
}

func (s *service) getAuthorizationFromQuery(r *http.Request) (string, error) {
	token := r.URL.Query().Get("access_token")

	if token == "" {
		return "", api.NewHttpError(http.StatusUnauthorized, "invalid token", nil)
	}

	return token, nil
}
