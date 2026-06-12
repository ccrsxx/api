package contents_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/contents"
	"github.com/ccrsxx/api/internal/test"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	db := &test.MockQuerier{
		ListContentByTypeFn: func(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, nil
		},
		UpsertContentFn: func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{}, nil
		},
	}

	service := contents.NewService(contents.ServiceConfig{Database: db})

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
