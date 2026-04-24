package jellyfin_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/jellyfin"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := jellyfin.NewService(jellyfin.ServiceConfig{})

	authService := auth.NewService(auth.ServiceConfig{})

	authMiddleware := auth.NewMiddleware(authService)

	jellyfin.LoadRoutes(jellyfin.Config{Router: mux, Service: service, AuthMiddleware: authMiddleware})

	tests := []test.RouteTestCase{

		{
			Path:   "/jellyfin/currently-playing",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
