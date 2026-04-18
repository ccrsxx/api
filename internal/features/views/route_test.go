package views

import (
	"context"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	db := newMockQuerier()

	svc := NewService(ServiceConfig{Database: db})

	ctx := context.Background()

	authMw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

	LoadRoutes(Config{
		Router:         mux,
		Service:        svc,
		AppContext:     ctx,
		AuthMiddleware: authMw,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/views/test-slug",
			Method: http.MethodGet,
		},
		{
			Path:   "/views/test-slug",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
