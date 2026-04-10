package sse

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
			Path:       "/sse",
			Method:     http.MethodGet,
			StatusCode: http.StatusUnauthorized,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
