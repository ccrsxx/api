package likes

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
	getContentBySlugFn     func(ctx context.Context, slug string) (sqlc.Content, error)
	getContentLikeStatusFn func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error)
	createContentLikeFn    func(ctx context.Context, arg sqlc.CreateContentLikeParams) error
	upsertIPAddressFn      func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error)
}

func (m *mockQuerier) GetContentBySlug(ctx context.Context, slug string) (sqlc.Content, error) {
	return m.getContentBySlugFn(ctx, slug)
}

func (m *mockQuerier) GetContentLikeStatus(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
	return m.getContentLikeStatusFn(ctx, arg)
}

func (m *mockQuerier) CreateContentLike(ctx context.Context, arg sqlc.CreateContentLikeParams) error {
	return m.createContentLikeFn(ctx, arg)
}

func (m *mockQuerier) UpsertIPAddress(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
	return m.upsertIPAddressFn(ctx, ipAddress)
}

var mockContentID = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
var mockIPAddressID = pgtype.UUID{Bytes: [16]byte{2}, Valid: true}

var mockContentLikeStatus = sqlc.GetContentLikeStatusRow{Likes: 10, UserLikes: 2}

func newMockQuerier() *mockQuerier {
	return &mockQuerier{
		getContentBySlugFn: func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{ID: mockContentID, Slug: slug}, nil
		},
		getContentLikeStatusFn: func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
			return mockContentLikeStatus, nil
		},
		createContentLikeFn: func(ctx context.Context, arg sqlc.CreateContentLikeParams) error {
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
			t.Fatalf("expected HTTPError, got %v", err)
		}

		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 HTTPError, got %v", err)
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

func TestService_getLikeStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		status, err := svc.getLikeStatus(context.Background(), "test-slug", mockIPAddressID)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if status.Likes != mockContentLikeStatus.Likes {
			t.Fatalf("got %d, want %d", status.Likes, mockContentLikeStatus.Likes)
		}

		if status.UserLikes != mockContentLikeStatus.UserLikes {
			t.Errorf("got %d, want %d", status.UserLikes, mockContentLikeStatus.UserLikes)
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentLikeStatusFn = func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
			return sqlc.GetContentLikeStatusRow{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.getLikeStatus(context.Background(), "test-slug", mockIPAddressID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get content like status error") {
			t.Errorf("got %v, want get content like status error", err)
		}
	})
}

func TestService_GetLikeStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		status, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if status.Likes != mockContentLikeStatus.Likes {
			t.Fatalf("got %d, want %d", status.Likes, mockContentLikeStatus.Likes)
		}

		if status.UserLikes != mockContentLikeStatus.UserLikes {
			t.Errorf("got %d, want %d", status.UserLikes, mockContentLikeStatus.UserLikes)
		}
	})

	t.Run("Get Content Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Upsert IP Error", func(t *testing.T) {
		db := newMockQuerier()

		db.upsertIPAddressFn = func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
			return sqlc.IpAddress{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "upsert ip address error") {
			t.Errorf("got %v, want upsert ip address error", err)
		}
	})
}

func TestService_IncrementLike(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		status, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if status.Likes != mockContentLikeStatus.Likes {
			t.Fatalf("got %d, want %d", status.Likes, mockContentLikeStatus.Likes)
		}

		if status.UserLikes != mockContentLikeStatus.UserLikes {
			t.Errorf("got %d, want %d", status.UserLikes, mockContentLikeStatus.UserLikes)
		}
	})

	t.Run("Get Content Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Upsert IP Error", func(t *testing.T) {
		db := newMockQuerier()

		db.upsertIPAddressFn = func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
			return sqlc.IpAddress{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "upsert ip address error") {
			t.Errorf("got %v, want upsert ip address error", err)
		}
	})

	t.Run("Get Initial Like Status Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentLikeStatusFn = func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
			return sqlc.GetContentLikeStatusRow{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Likes Limit Reached", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentLikeStatusFn = func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
			return sqlc.GetContentLikeStatusRow{Likes: 10, UserLikes: 5}, nil
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("expected HTTPError, got %v", err)
		}

		if httpErr.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403 HTTPError, got %v", err)
		}
	})

	t.Run("Create Content Like Error", func(t *testing.T) {
		db := newMockQuerier()

		db.createContentLikeFn = func(ctx context.Context, arg sqlc.CreateContentLikeParams) error {
			return errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create content like error") {
			t.Errorf("got %v, want create content like error", err)
		}
	})
}
