package contents_test

import (
	"context"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/contents"
	"github.com/ccrsxx/api/internal/test"
)

func newMockQuerier() *test.MockQuerier {
	return &test.MockQuerier{
		ListContentByTypeFn: func(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error) {
			return []sqlc.ListContentByTypeRow{
				{Slug: "test-post", Type: "blog", Views: 10, Likes: 5},
				{Slug: "another-post", Type: "blog", Views: 20, Likes: 8},
			}, nil
		},
		UpsertContentFn: func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{Slug: arg.Slug, Type: arg.Type}, nil
		},
	}
}

func TestService_GetContentsData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		data, err := svc.GetContentsData(context.Background(), "blog")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if len(data) != 2 {
			t.Fatalf("got %d, want 2", len(data))
		}

		if data[0].Type != "blog" {
			t.Fatalf("got %s, want blog", data[0].Type)
		}
	})

	t.Run("Valid Empty Data", func(t *testing.T) {
		db := newMockQuerier()

		db.ListContentByTypeFn = func(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, nil
		}

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		data, err := svc.GetContentsData(context.Background(), "blog")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if data == nil {
			t.Fatal("got nil, want empty slice")
		}

		if len(data) != 0 {
			t.Fatalf("got %d, want 0", len(data))
		}
	})

	t.Run("Invalid Content Type", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		_, err := svc.GetContentsData(context.Background(), "invalid")

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.ListContentByTypeFn = func(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, context.DeadlineExceeded
		}

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		_, err := svc.GetContentsData(context.Background(), "blog")

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})
}

func TestService_UpsertContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})

		input := contents.UpsertContentInput{
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

		if content.Type != "blog" {
			t.Fatalf("got %s, want blog", content.Type)
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.UpsertContentFn = func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{}, context.DeadlineExceeded
		}

		svc := contents.NewService(contents.ServiceConfig{Database: db})

		input := contents.UpsertContentInput{
			Slug: "new-post",
			Type: "blog",
		}

		_, err := svc.UpsertContent(context.Background(), input)

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})
}
