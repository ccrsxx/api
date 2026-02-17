package sse

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/src/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	LoadRoutes(mux)

	tests := []test.RouteTestCase{
		{
			Path:       "/sse",
			Method:     http.MethodGet,
			StatusCode: http.StatusUnauthorized,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
