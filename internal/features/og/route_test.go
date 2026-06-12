package og_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/og"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := og.NewService(og.ServiceConfig{})

	og.LoadRoutes(og.Config{
		Router:  mux,
		Service: service,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/og",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
