package pixiv_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/pixiv"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	mock := &mockPixivClient{
		imageErr: errors.New("not found"),
	}

	service := pixiv.NewService(pixiv.ServiceConfig{
		Client: mock,
	})

	authService := auth.NewService(auth.ServiceConfig{})

	authMiddleware := auth.NewMiddleware(authService)

	pixiv.LoadRoutes(pixiv.Config{Router: mux, Service: service, AuthMiddleware: authMiddleware})

	tests := []test.RouteTestCase{
		{
			Path:   "/pixiv/bookmarks",
			Method: http.MethodGet,
		},
		{
			Path:   "/pixiv/bookmarks/all",
			Method: http.MethodGet,
		},
		{
			Path:   "/pixiv/image/{url}",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
