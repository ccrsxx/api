package sse

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := NewService(Config{})

	LoadRoutes(mux, service)

	tests := []test.RouteTestCase{
		{
			Path:       "/sse",
			Method:     http.MethodGet,
			StatusCode: http.StatusUnauthorized,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
