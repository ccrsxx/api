package og

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	service := NewService(ServiceConfig{})

	LoadRoutes(mux, service, Config{})

	tests := []test.RouteTestCase{
		{
			Path:   "/og",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
