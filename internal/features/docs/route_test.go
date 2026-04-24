package docs_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/docs"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	docs.LoadRoutes(docs.Config{Router: mux})

	tests := []test.RouteTestCase{
		{
			Path:   "/docs",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
