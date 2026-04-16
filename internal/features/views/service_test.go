package views

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockQuerier struct {
	getContentBySlugFn    func(ctx context.Context, slug string) (sqlc.Content, error)
	getContentViewCountFn func(ctx context.Context, slug string) (int32, error)
	createContentViewFn   func(ctx context.Context, arg sqlc.CreateContentViewParams) error
	upsertIPAddressFn     func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error)
}

func (m *mockQuerier) GetContentBySlug(ctx context.Context, slug string) (sqlc.Content, error) {
	return m.getContentBySlugFn(ctx, slug)
}

func (m *mockQuerier) GetContentViewCount(ctx context.Context, slug string) (int32, error) {
	return m.getContentViewCountFn(ctx, slug)
}

func (m *mockQuerier) CreateContentView(ctx context.Context, arg sqlc.CreateContentViewParams) error {
	return m.createContentViewFn(ctx, arg)
}

func (m *mockQuerier) UpsertIPAddress(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
	return m.upsertIPAddressFn(ctx, ipAddress)
}

var mockContentID = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
var mockIPAddressID = pgtype.UUID{Bytes: [16]byte{2}, Valid: true}

var mockContentViewCount = int32(5)

func newMockQuerier() *mockQuerier {
	return &mockQuerier{
		getContentBySlugFn: func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{ID: mockContentID, Slug: slug}, nil
		},
		getContentViewCountFn: func(ctx context.Context, slug string) (int32, error) {
			return mockContentViewCount, nil
		},
		createContentViewFn: func(ctx context.Context, arg sqlc.CreateContentViewParams) error {
			return nil
		},
		upsertIPAddressFn: func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
			return sqlc.IpAddress{ID: mockIPAddressID, IpAddress: ipAddress}, nil
		},
	}
}

func TestService_getContentBySlug(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		content, err := svc.getContentBySlug(context.Background(), "test-slug")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if content.Slug != "test-slug" {
			t.Errorf("got %s, want test-slug", content.Slug)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, pgx.ErrNoRows
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.getContentBySlug(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("got status %d, want 404", httpErr.StatusCode)
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.getContentBySlug(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get content by slug error") {
			t.Errorf("got %v, want get content by slug error", err)
		}
	})
}

func TestService_GetViewCount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		views, err := svc.GetViewCount(context.Background(), "test-slug")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if views.Views != mockContentViewCount {
			t.Errorf("got %d, want %d", views.Views, mockContentViewCount)
		}
	})

	t.Run("Get Content Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.GetViewCount(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Get View Count Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentViewCountFn = func(ctx context.Context, slug string) (int32, error) {
			return 0, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.GetViewCount(context.Background(), "test-slug")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get view count error") {
			t.Errorf("got %v, want get view count error", err)
		}
	})
}

func TestService_IncrementView(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		views, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if views.Views != mockContentViewCount {
			t.Errorf("got %d, want %d", views.Views, mockContentViewCount)
		}
	})

	t.Run("Get Content Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Upsert IP Error", func(t *testing.T) {
		db := newMockQuerier()

		db.upsertIPAddressFn = func(ctx context.Context, ip string) (sqlc.IpAddress, error) {
			return sqlc.IpAddress{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

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

		db.createContentViewFn = func(ctx context.Context, arg sqlc.CreateContentViewParams) error {
			return errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

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

		db.getContentViewCountFn = func(ctx context.Context, slug string) (int32, error) {
			return 0, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementView(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get view count after increment error") {
			t.Errorf("got %v, want get view count after increment error", err)
		}
	})
}
