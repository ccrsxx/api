package contents_test

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
	"github.com/ccrsxx/api/internal/features/contents"
	"github.com/ccrsxx/api/internal/test"
)

var validPath = "/?type=blog"

func TestController_GetContentsData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentsData(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[[]sqlc.ListContentByTypeRow]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		t.Logf("response: %+v", res)

		if len(res.Data) != 2 {
			t.Fatalf("got %d, want 2", len(res.Data))
		}

		if res.Data[0].Type != "blog" {
			t.Fatalf("got %s, want blog", res.Data[0].Type)
		}

	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.ListContentByTypeFn = func(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error) {
			return nil, errors.New("db error")
		}

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.GetContentsData(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, validPath, nil)

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetContentsData(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_UpsertContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		input := contents.UpsertContentInput{Slug: "new-post", Type: "blog"}

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

		if res.Data.Type != "blog" {
			t.Fatalf("got %s, want blog", res.Data.Type)
		}
	})

	t.Run("Decode Error", func(t *testing.T) {
		db := newMockQuerier()

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		r := httptest.NewRequest(http.MethodPost, validPath, nil)

		w := httptest.NewRecorder()

		ctrl.UpsertContent(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", w.Code)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := newMockQuerier()

		db.UpsertContentFn = func(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error) {
			return sqlc.Content{}, errors.New("db error")
		}

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		input := contents.UpsertContentInput{Slug: "new-post", Type: "blog"}

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

		svc := contents.NewService(contents.ServiceConfig{Database: db})
		ctrl := contents.NewController(svc)

		input := contents.UpsertContentInput{Slug: "new-post", Type: "blog"}

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
