package likes_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/likes"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	db := &test.MockQuerier{
		GetContentBySlugFn:     mockGetContentBySlugFn,
		GetContentLikeStatusFn: mockGetContentLikeStatusFn,
		UpsertIPAddressFn:      mockUpsertIPAddressFn,
	}

	svc := likes.NewService(likes.ServiceConfig{Database: db})

	authMw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

	likes.LoadRoutes(likes.Config{
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
