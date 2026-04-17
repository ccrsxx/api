package statistics

import (
	"context"
	"errors"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
)

type mockQuerier struct {
	getContentStatsByTypeFn func(ctx context.Context, kind string) (sqlc.GetContentStatsByTypeRow, error)
}

func (m *mockQuerier) GetContentStatsByType(ctx context.Context, kind string) (sqlc.GetContentStatsByTypeRow, error) {
	return m.getContentStatsByTypeFn(ctx, kind)
}

var mockContentStatsByType = sqlc.GetContentStatsByTypeRow{
	TotalPosts: 5,
	TotalViews: 10,
	TotalLikes: 15,
}

func newMockQuerier() *mockQuerier {
	return &mockQuerier{
		getContentStatsByTypeFn: func(ctx context.Context, kind string) (sqlc.GetContentStatsByTypeRow, error) {
			return mockContentStatsByType, nil
		},
	}
}

func TestService_GetContentStatistics(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		stats, err := svc.GetContentStatistics(context.Background(), "blog")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if stats.Type != "blog" {
			t.Fatalf("got %s, want blog", stats.Type)
		}

		if stats.TotalPosts != mockContentStatsByType.TotalPosts {
			t.Fatalf("got %d, want %d", stats.TotalPosts, mockContentStatsByType.TotalPosts)
		}

		if stats.TotalViews != mockContentStatsByType.TotalViews {
			t.Fatalf("got %d, want %d", stats.TotalViews, mockContentStatsByType.TotalViews)
		}

		if stats.TotalLikes != mockContentStatsByType.TotalLikes {
			t.Errorf("got %d, want %d", stats.TotalLikes, mockContentStatsByType.TotalLikes)
		}
	})

	t.Run("Invalid Content Type", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		_, err := svc.GetContentStatistics(context.Background(), "invalid")

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentStatsByTypeFn = func(ctx context.Context, kind string) (sqlc.GetContentStatsByTypeRow, error) {
			return sqlc.GetContentStatsByTypeRow{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})

		_, err := svc.GetContentStatistics(context.Background(), "blog")

		if err == nil {
			t.Error("got nil, want error")
		}
	})
}
