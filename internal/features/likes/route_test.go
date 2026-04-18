package likes

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	svc := NewService(ServiceConfig{Database: newMockQuerier()})

	authMw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

	LoadRoutes(Config{
		Router:         mux,
		Service:        svc,
		AuthMiddleware: authMw,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/likes/test-slug",
			Method: http.MethodGet,
		},
		{
			Path:   "/likes/test-slug",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
