package auth_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
	"golang.org/x/oauth2"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	svc := auth.NewService(auth.ServiceConfig{
		GithubOauthConfig: &oauth2.Config{},
	})

	authMw := auth.NewMiddleware(svc)

	auth.LoadRoutes(auth.Config{
		Router:         mux,
		Service:        svc,
		AuthMiddleware: authMw,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/auth/me",
			Method: http.MethodGet,
		},
		{
			Path:   "/auth/github/login",
			Method: http.MethodGet,
		},
		{
			Path:   "/auth/github/logout",
			Method: http.MethodPost,
		},
		{
			Path:   "/auth/github/callback",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
