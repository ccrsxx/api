package home

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
			Path:   "/",
			Method: http.MethodGet,
		},
		{
			Path:   "/",
			Host:   "ip.example.com",
			Method: http.MethodGet,
		},
		{
			Path:   "/",
			Host:   "ipinfo.example.com",
			Method: http.MethodGet,
		},
		{
			Path:   "/",
			Host:   "headers.example.com",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
