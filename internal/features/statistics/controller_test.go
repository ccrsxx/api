package statistics_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/statistics"
	"github.com/ccrsxx/api/internal/test"
)

var validPath = "/?type=blog"

func TestController_GetContentsStatistics(t *testing.T) {
	mockStatsFn := func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
		return mockContentStatsByType, nil
	}

	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{GetContentStatsByTypeFn: mockStatsFn}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})
		ctrl := statistics.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentsStatistics(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[statistics.ContentsStatistics]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Type != "blog" {
			t.Fatalf("got %s, want %s", res.Data.Type, "blog")
		}

		if res.Data.TotalPosts != mockContentStatsByType.TotalPosts {
			t.Fatalf("got %d, want %d", res.Data.TotalPosts, mockContentStatsByType.TotalPosts)
		}

		if res.Data.TotalViews != mockContentStatsByType.TotalViews {
			t.Fatalf("got %d, want %d", res.Data.TotalViews, mockContentStatsByType.TotalViews)
		}

		if res.Data.TotalLikes != mockContentStatsByType.TotalLikes {
			t.Errorf("got %d, want %d", res.Data.TotalLikes, mockContentStatsByType.TotalLikes)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetContentStatsByTypeFn: func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
				return sqlc.GetContentStatsByTypeRow{}, errors.New("db error")
			},
		}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})
		ctrl := statistics.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentsStatistics(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := &test.MockQuerier{GetContentStatsByTypeFn: mockStatsFn}

		svc := statistics.NewService(statistics.ServiceConfig{Database: db})
		ctrl := statistics.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetContentsStatistics(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
