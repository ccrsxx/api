package navidrome_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/auth"
	navidromeFeature "github.com/ccrsxx/api/internal/features/navidrome"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	mock := &mockNavidromeClient{
		coverArtErr: errors.New("not found"),
	}

	service := navidromeFeature.NewService(navidromeFeature.ServiceConfig{
		Client: mock,
	})

	authService := auth.NewService(auth.ServiceConfig{})

	authMiddleware := auth.NewMiddleware(authService)

	navidromeFeature.LoadRoutes(navidromeFeature.Config{
		Router:         mux,
		Service:        service,
		AuthMiddleware: authMiddleware,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/navidrome/currently-playing",
			Method: http.MethodGet,
		},
		{
			Path:   "/navidrome/cover-art/test-id",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
