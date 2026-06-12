package favicon_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/favicon"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	favicon.LoadRoutes(favicon.Config{Router: mux})

	tests := []test.RouteTestCase{
		{
			Path:   "/favicon.ico",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
