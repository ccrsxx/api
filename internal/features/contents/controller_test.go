package contents

import (
	"bytes"
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

func TestController_GetContentData(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentData(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[ContentData]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Type != "blog" {
			t.Fatalf("got %s, want blog", res.Data.Type)
		}

		if len(res.Data.Data) != 2 {
			t.Fatalf("got %d, want 2", len(res.Data.Data))
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.listContentByTypeFn = func(ctx context.Context, kind string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentData(w, r)

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

		ctrl.GetContentData(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_UpsertContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		input := UpsertContentInput{Slug: "new-post", Type: "blog"}

		jsonBytes, err := json.Marshal(input)

		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, validPath, bytes.NewReader(jsonBytes))

		w := httptest.NewRecorder()

		ctrl.UpsertContent(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", w.Code)
		}

		var res api.SuccessResponse[sqlc.Content]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.Slug != "new-post" {
			t.Fatalf("got %s, want new-post", res.Data.Slug)
		}

		if res.Data.Kind != "blog" {
			t.Fatalf("got %s, want blog", res.Data.Kind)
		}
	})

	t.Run("Decode Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		r := httptest.NewRequest(http.MethodPost, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.UpsertContent(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", w.Code)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.upsertContentFn = func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		input := UpsertContentInput{Slug: "new-post", Type: "blog"}

		jsonBytes, err := json.Marshal(input)

		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, validPath, bytes.NewReader(jsonBytes))

		w := httptest.NewRecorder()

		ctrl.UpsertContent(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := NewService(ServiceConfig{Database: db})
		ctrl := NewController(svc)

		input := UpsertContentInput{Slug: "new-post", Type: "blog"}

		jsonBytes, err := json.Marshal(input)

		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, validPath, bytes.NewReader(jsonBytes))

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.UpsertContent(errWriter, r)

		if w.Code != http.StatusCreated {
			t.Errorf("got %d, want %d", w.Code, http.StatusCreated)
		}
	})

}
