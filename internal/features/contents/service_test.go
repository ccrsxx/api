package contents

import (
	"context"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
)

type mockQuerier struct {
	listContentByTypeFn func(ctx context.Context, kind string) ([]sqlc.ListContentByTypeRow, error)
	upsertContentFn     func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error)
}

func (m *mockQuerier) ListContentByType(ctx context.Context, kind string) ([]sqlc.ListContentByTypeRow, error) {
	return m.listContentByTypeFn(ctx, kind)
}

func (m *mockQuerier) UpsertContent(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
	return m.upsertContentFn(ctx, arg)
}

func newMockQuerier() *mockQuerier {
	return &mockQuerier{
		listContentByTypeFn: func(ctx context.Context, kind string) ([]sqlc.ListContentByTypeRow, error) {
			return []sqlc.ListContentByTypeRow{
				{Slug: "test-post", Views: 10, Likes: 5},
				{Slug: "another-post", Views: 20, Likes: 8},
			}, nil
		},
		upsertContentFn: func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{Slug: arg.Slug, Kind: arg.Kind}, nil
		},
	}
}

func TestService_GetContentData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		data, err := svc.GetContentData(context.Background(), "blog")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if data.Type != "blog" {
			t.Fatalf("got %s, want blog", data.Type)
		}

		if len(data.Data) != 2 {
			t.Fatalf("got %d, want 2", len(data.Data))
		}
	})

	t.Run("Valid Empty Data", func(t *testing.T) {
		db := newMockQuerier()

		db.listContentByTypeFn = func(ctx context.Context, kind string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, nil
		}

		svc := NewService(ServiceConfig{Database: db})
		data, err := svc.GetContentData(context.Background(), "blog")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if data.Data == nil {
			t.Fatal("got nil, want empty slice")
		}

		if len(data.Data) != 0 {
			t.Fatalf("got %d, want 0", len(data.Data))
		}
	})

	t.Run("Invalid Content Type", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		_, err := svc.GetContentData(context.Background(), "invalid")

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.listContentByTypeFn = func(ctx context.Context, kind string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, context.DeadlineExceeded
		}

		svc := NewService(ServiceConfig{Database: db})
		_, err := svc.GetContentData(context.Background(), "blog")

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})
}

func TestService_UpsertContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		input := UpsertContentInput{
			Slug: "new-post",
			Type: "blog",
		}

		content, err := svc.UpsertContent(context.Background(), input)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if content.Slug != "new-post" {
			t.Fatalf("got %s, want new-post", content.Slug)
		}

		if content.Kind != "blog" {
			t.Fatalf("got %s, want blog", content.Kind)
		}
	})

	t.Run("Invalid Content Type", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})

		input := UpsertContentInput{
			Slug: "new-post",
			Type: "invalid",
		}

		_, err := svc.UpsertContent(context.Background(), input)

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.upsertContentFn = func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{}, context.DeadlineExceeded
		}

		svc := NewService(ServiceConfig{Database: db})

		input := UpsertContentInput{
			Slug: "new-post",
			Type: "blog",
		}

		_, err := svc.UpsertContent(context.Background(), input)

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})
}
