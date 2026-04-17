package views

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
	"github.com/jackc/pgx/v5/pgtype"
)

func TestController_GetViewCount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.GetViewCount(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[ViewCount]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Views != int64(mockContentViewCount) {
			t.Errorf("got %d, want %d", res.Data.Views, mockContentViewCount)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getTotalContentMeta = func(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error) {
			return sqlc.GetTotalContentMetaRow{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.GetViewCount(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/test-slug", nil)

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetViewCount(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_IncrementView(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodPost, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.IncrementView(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[ViewCount]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Views != int64(mockContentViewCount) {
			t.Errorf("got %d, want %d", res.Data.Views, mockContentViewCount)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentBySlugFn = func(ctx context.Context, slug string) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("not found")
		}

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodPost, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.IncrementView(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodPost, "/test-slug", nil)

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.IncrementView(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
