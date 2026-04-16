package likes

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

func TestController_GetLikeStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.GetLikeStatus(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[sqlc.GetContentLikeStatusRow]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Likes != mockContentLikeStatus.Likes {
			t.Fatalf("got %d, want %d", res.Data.Likes, mockContentLikeStatus.Likes)
		}

		if res.Data.UserLikes != mockContentLikeStatus.UserLikes {
			t.Errorf("got %d, want %d", res.Data.Likes, mockContentLikeStatus.UserLikes)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.getContentLikeStatusFn = func(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error) {
			return sqlc.GetContentLikeStatusRow{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.GetLikeStatus(w, r)

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

		ctrl.GetLikeStatus(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_IncrementLike(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodPost, "/test-slug", nil)

		w := httptest.NewRecorder()

		ctrl.IncrementLike(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 201", w.Code)
		}

		var res api.SuccessResponse[sqlc.GetContentLikeStatusRow]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Likes != mockContentLikeStatus.Likes {
			t.Fatalf("got %d, want %d", res.Data.Likes, mockContentLikeStatus.Likes)
		}

		if res.Data.UserLikes != mockContentLikeStatus.UserLikes {
			t.Errorf("got %d, want %d", res.Data.Likes, mockContentLikeStatus.UserLikes)
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

		ctrl.IncrementLike(w, r)

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

		ctrl.IncrementLike(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
