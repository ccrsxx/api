package statistics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/statistics"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	db := &test.MockQuerier{
		GetContentStatsByTypeFn: func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
			return sqlc.GetContentStatsByTypeRow{}, nil
		},
	}

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
