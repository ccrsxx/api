package contents_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/contents"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := contents.NewService(contents.ServiceConfig{Database: newMockQuerier()})

	contents.LoadRoutes(contents.Config{Router: mux, Service: service})

	tests := []test.RouteTestCase{
		{
			Path:   "/contents",
			Method: http.MethodGet,
		},
		{
			Path:   "/contents",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
