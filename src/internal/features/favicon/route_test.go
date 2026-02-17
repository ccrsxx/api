package favicon

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
			Path:   "/favicon.ico",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
