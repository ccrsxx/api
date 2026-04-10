package jellyfin

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := NewService(ServiceConfig{})

	authService := auth.NewService(auth.ServiceConfig{})

	authMiddleware := auth.NewMiddleware(authService)

	LoadRoutes(mux, service, authMiddleware)

	tests := []test.RouteTestCase{

		{
			Path:   "/jellyfin/currently-playing",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
