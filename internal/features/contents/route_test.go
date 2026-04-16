package contents

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
			Path:   "/content/blog",
			Method: http.MethodGet,
		},
		{
			Path:   "/content/project",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
