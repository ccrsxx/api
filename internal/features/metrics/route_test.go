package metrics_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/metrics"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	metrics.LoadRoutes(metrics.Config{Router: mux})

	tests := []test.RouteTestCase{
		{
			Path:   "/metrics",
			Method: http.MethodGet,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
