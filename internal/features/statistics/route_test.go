package statistics_test

import (
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/features/statistics"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	db := newMockQuerier()

	svc := statistics.NewService(statistics.ServiceConfig{Database: db})

	statistics.LoadRoutes(statistics.Config{Router: mux, Service: svc})

	tests := []test.RouteTestCase{
		{
			Path:   "/statistics",
			Method: http.MethodGet,
		},
		{
			Path:   "/statistics",
			Method: http.MethodPost,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
