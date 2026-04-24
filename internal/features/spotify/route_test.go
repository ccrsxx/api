package spotify_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/spotify"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := spotify.NewService(spotify.ServiceConfig{})

	authService := auth.NewService(auth.ServiceConfig{})

	authMiddleware := auth.NewMiddleware(authService)

	spotify.LoadRoutes(spotify.Config{Router: mux, Service: service, AuthMiddleware: authMiddleware})

	tests := []test.RouteTestCase{
		{
			Path:   "/spotify/currently-playing",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
