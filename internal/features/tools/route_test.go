package tools

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	LoadRoutes(mux)

	tests := []test.RouteTestCase{
		{
			Path:   "/tools/ip",
			Method: http.MethodGet,
		},
		{
			Path:   "/tools/headers",
			Method: http.MethodGet,
		},
		{
			Path:   "/tools/ipinfo",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
