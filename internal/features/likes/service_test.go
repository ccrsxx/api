package likes_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/likes"
	"github.com/ccrsxx/api/internal/test"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	likesCount = 10
	userLikes  = 1

	mockContentID   = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockIPAddressID = pgtype.UUID{Bytes: [16]byte{2}, Valid: true}

	mockContentLikeStatus = sqlc.GetContentLikeStatusRow{Likes: int64(likesCount), UserLikes: int64(userLikes)}

	mockIncrementContentLike = sqlc.IncrementContentLikeRow{Likes: int32(likesCount)}

	mockGetContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
		return sqlc.Content{ID: mockContentID, Slug: slug}, nil
	}

	mockGetContentLikeStatusFn = func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
		return mockContentLikeStatus, nil
	}

	mockIncrementContentLikeFn = func(ctx context.Context, arg sqlc.IncrementContentLikeParams) (sqlc.IncrementContentLikeRow, error) {
		return mockIncrementContentLike, nil
	}

	mockUpsertIPAddressFn = func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
		return sqlc.IpAddress{ID: mockIPAddressID, IpAddress: ipAddress}, nil
	}
)

func TestService_GetLikeStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn:     mockGetContentBySlugFn,
			GetContentLikeStatusFn: mockGetContentLikeStatusFn,
			UpsertIPAddressFn:      mockUpsertIPAddressFn,
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

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

	t.Run("Get Content Error (coverage)", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: func(ctx context.Context, slug string) (sqlc.Content, error) {
				return sqlc.Content{}, errors.New("db error")
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Upsert IP Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: mockGetContentBySlugFn,
			UpsertIPAddressFn: func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
				return sqlc.IpAddress{}, errors.New("db error")
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "upsert ip address error") {
			t.Errorf("got %v, want upsert ip address error", err)
		}
	})

	t.Run("Content Not Found (coverage)", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: func(ctx context.Context, slug string) (sqlc.Content, error) {
				return sqlc.Content{}, pgx.ErrNoRows
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

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

	t.Run("Get Like Status Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: mockGetContentBySlugFn,
			UpsertIPAddressFn:  mockUpsertIPAddressFn,
			GetContentLikeStatusFn: func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
				return sqlc.GetContentLikeStatusRow{}, errors.New("db error")
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.GetLikeStatus(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get content like status error") {
			t.Errorf("got %v, want get content like status error", err)
		}
	})
}

func TestService_IncrementLike(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn:     mockGetContentBySlugFn,
			GetContentLikeStatusFn: mockGetContentLikeStatusFn,
			IncrementContentLikeFn: mockIncrementContentLikeFn,
			UpsertIPAddressFn:      mockUpsertIPAddressFn,
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

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

	t.Run("Upsert IP Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: mockGetContentBySlugFn,
			UpsertIPAddressFn: func(ctx context.Context, ipAddress string) (sqlc.IpAddress, error) {
				return sqlc.IpAddress{}, errors.New("db error")
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "upsert ip address error") {
			t.Errorf("got %v, want upsert ip address error", err)
		}
	})

	t.Run("Get Initial Like Status Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: mockGetContentBySlugFn,
			UpsertIPAddressFn:  mockUpsertIPAddressFn,
			GetContentLikeStatusFn: func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
				return sqlc.GetContentLikeStatusRow{}, errors.New("db error")
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Likes Limit Reached", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn: mockGetContentBySlugFn,
			UpsertIPAddressFn:  mockUpsertIPAddressFn,
			GetContentLikeStatusFn: func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
				return sqlc.GetContentLikeStatusRow{Likes: 10, UserLikes: 5}, nil
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("expected HTTPError, got %v", err)
		}

		if httpErr.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected 422 HTTPError, got %v", err)
		}
	})

	t.Run("Create Content Like Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentBySlugFn:     mockGetContentBySlugFn,
			GetContentLikeStatusFn: mockGetContentLikeStatusFn,
			UpsertIPAddressFn:      mockUpsertIPAddressFn,
			IncrementContentLikeFn: func(ctx context.Context, arg sqlc.IncrementContentLikeParams) (sqlc.IncrementContentLikeRow, error) {
				return sqlc.IncrementContentLikeRow{}, errors.New("db error")
			},
		}

		svc := likes.NewService(likes.ServiceConfig{Database: db})

		_, err := svc.IncrementLike(context.Background(), "test-slug", "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create content like error") {
			t.Errorf("got %v, want create content like error", err)
		}
	})
}
