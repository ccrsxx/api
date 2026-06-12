package guestbook_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/guestbook"
	"github.com/ccrsxx/api/internal/test"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestLoadRoutes(t *testing.T) {
	mux := http.NewServeMux()

	db := &test.MockQuerier{
		ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
			return []sqlc.ListGuestbookRow{}, nil
		},
		CreateGuestbookFn: func(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error) {
			return sqlc.CreateGuestbookRow{}, nil
		},
		DeleteGuestbookFn: func(ctx context.Context, id pgtype.UUID) error {
			return nil
		},
		GetGuestbookByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error) {
			return sqlc.Guestbook{ID: id}, nil
		},
	}

	svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

	authMw := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{}))

	guestbook.LoadRoutes(guestbook.Config{
		Router:         mux,
		Service:        svc,
		AuthMiddleware: authMw,
	})

	tests := []test.RouteTestCase{
		{
			Path:   "/guestbook/",
			Method: http.MethodGet,
		},
		{
			Path:   "/guestbook/",
			Method: http.MethodPost,
		},
		{
			Path:   "/guestbook/123",
			Method: http.MethodDelete,
		},
	}

	test.AssertRoutes(t, mux, tests)
}
