package likes

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := NewService(ServiceConfig{Database: newMockQuerier()})

	LoadRoutes(Config{Router: mux, Service: service})

	tests := []test.RouteTestCase{
		{
			Path:   "/likes/test-slug",
			Method: http.MethodGet,
		},
		{
			Path:   "/likes/test-slug",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
