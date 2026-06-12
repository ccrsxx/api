package statistics_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/statistics"
	"github.com/ccrsxx/api/internal/test"
)

var mockContentStatsByType = sqlc.GetContentStatsByTypeRow{
	TotalPosts: 5,
	TotalViews: 10,
	TotalLikes: 15,
}

func TestService_GetContentsStatistics(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentStatsByTypeFn: func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
				return mockContentStatsByType, nil
			},
		}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})
		stats, err := svc.GetContentsStatistics(context.Background(), "blog")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
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

	t.Run("Success Empty Type", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentStatsByTypeFn: func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
				return mockContentStatsByType, nil
			},
		}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})
		stats, err := svc.GetContentsStatistics(context.Background(), "")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if stats.Type != "all" {
			t.Errorf("got %s, want all", stats.Type)
		}
	})

	t.Run("Invalid Content Type", func(t *testing.T) {
		db := &test.MockQuerier{}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})
		_, err := svc.GetContentsStatistics(context.Background(), "invalid")

		if err == nil {
			t.Fatal("got nil, want error")
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentStatsByTypeFn: func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
				return sqlc.GetContentStatsByTypeRow{}, errors.New("db error")
			},
		}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})

		_, err := svc.GetContentsStatistics(context.Background(), "blog")

		if err == nil {
			t.Error("got nil, want error")
		}
	})
}
