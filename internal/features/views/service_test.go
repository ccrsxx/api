package views_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/views"
	"github.com/ccrsxx/api/internal/test"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	mockContentID   = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockIPAddressID = pgtype.UUID{Bytes: [16]byte{2}, Valid: true}

	mockContentViewCount = 5
	mockContentLikeCount = 5

	mockIncrementContentView = sqlc.IncrementContentViewRow{
		Views: int32(mockContentViewCount),
		Likes: int32(mockContentLikeCount),
	}

	mockTotalContentMeta = sqlc.GetTotalContentMetaRow{
		TotalViews: int64(mockContentViewCount),
		TotalLikes: int64(mockContentLikeCount),
	}
)

func newMockQuerier() *test.MockQuerier {
	return &test.MockQuerier{
		GetContentBySlugFn: func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{ID: mockContentID, Slug: slug}, nil
		},
		GetTotalContentMetaFn: func(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error) {
			return mockTotalContentMeta, nil
		},
		IncrementContentViewFn: func(ctx context.Context, arg sqlc.IncrementContentViewParams) (sqlc.IncrementContentViewRow, error) {
			return mockIncrementContentView, nil
		},
		UpsertIPAddressFn: func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
			return sqlc.IpAddress{ID: mockIPAddressID, IpAddress: ipAddress}, nil
		},
	}
}

func TestService_GetViewCount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := views.NewService(views.ServiceConfig{Database: db})

		viewCount, err := svc.GetViewCount(context.Background(), "test-slug")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if viewCount.Views != int64(mockContentViewCount) {
			t.Errorf("got %d, want %d", viewCount.Views, mockContentViewCount)
		}
	})

	t.Run("Get Content Error", func(t *testing.T) {
		db := newMockQuerier()

		db.GetContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.GetViewCount(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Get View Count Error", func(t *testing.T) {
		db := newMockQuerier()

		db.GetTotalContentMetaFn = func(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error) {
			return sqlc.GetTotalContentMetaRow{}, errors.New("db error")
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.GetViewCount(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get view count error") {
			t.Errorf("got %v, want get view count error", err)
		}
	})

	t.Run("Content Not Found (coverage)", func(t *testing.T) {
		db := newMockQuerier()

		db.GetContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, pgx.ErrNoRows
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.GetViewCount(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want *api.HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("got %d, want 404", httpErr.StatusCode)
		}
	})
}

func TestService_IncrementView(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := views.NewService(views.ServiceConfig{Database: db})

		viewCount, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if viewCount.Views != int64(mockContentViewCount) {
			t.Errorf("got %d, want %d", viewCount.Views, mockContentViewCount)
		}
	})

	t.Run("Get Content Error", func(t *testing.T) {
		db := newMockQuerier()

		db.GetContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Upsert IP Error", func(t *testing.T) {
		db := newMockQuerier()

		db.UpsertIPAddressFn = func(ctx context.Context, ip string) (sqlc.IpAddress, error) {
			return sqlc.IpAddress{}, errors.New("db error")
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "upsert ip address error") {
			t.Errorf("got %v, want upsert ip address error", err)
		}
	})

	t.Run("Create Content View Error", func(t *testing.T) {
		db := newMockQuerier()

		db.IncrementContentViewFn = func(ctx context.Context, arg sqlc.IncrementContentViewParams) (sqlc.IncrementContentViewRow, error) {
			return sqlc.IncrementContentViewRow{}, errors.New("db error")
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create content view error") {
			t.Errorf("got %v, want create content view error", err)
		}
	})

	t.Run("Get View Count Error", func(t *testing.T) {
		db := newMockQuerier()

		db.GetTotalContentMetaFn = func(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error) {
			return sqlc.GetTotalContentMetaRow{}, errors.New("db error")
		}

		svc := views.NewService(views.ServiceConfig{Database: db})

		_, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get view count after increment error") {
			t.Errorf("got %v, want get view count after increment error", err)
		}
	})
}
