package statistics

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/test"
)

var validPath = "/?type=blog"

func TestController_GetContentsStatistics(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentsStatistics(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[sqlc.GetContentStatsByTypeRow]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
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
		db := newMockQuerier()

		db.getContentStatsByTypeFn = func(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error) {
			return sqlc.GetContentStatsByTypeRow{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentsStatistics(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetContentsStatistics(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
