package guestbook_test

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
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/guestbook"
	"github.com/ccrsxx/api/internal/test"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestController_GetGuestbook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
				return mockListGuestbookRow, nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})
		ctrl := guestbook.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		ctrl.GetGuestbook(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", w.Code)
		}

		var res api.SuccessResponse[[]sqlc.ListGuestbookRow]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(res.Data) != 1 || res.Data[0].ID != mockGuestbookID {
			t.Errorf("got invalid response data: %+v", res.Data)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := &test.MockQuerier{
			ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
				return nil, errors.New("db error")
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})
		ctrl := guestbook.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		ctrl.GetGuestbook(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := &test.MockQuerier{
			ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
				return mockListGuestbookRow, nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})
		ctrl := guestbook.NewController(svc)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.GetGuestbook(errWriter, r)

		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestController_CreateGuestbook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateGuestbookFn: func(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error) {
				return mockCreateGuestbookRow, nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{
			Database:    db,
			EmailClient: &mockEmailClient{},
		})

		ctrl := guestbook.NewController(svc)

		body := []byte(`{"text":"Hello world!"}`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		ctx := auth.SetUserContext(r.Context(), mockUser)
		r = r.WithContext(ctx)

		w := httptest.NewRecorder()

		ctrl.CreateGuestbook(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", w.Code)
		}

		var res api.SuccessResponse[sqlc.CreateGuestbookRow]

		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Data.ID != mockGuestbookID {
			t.Errorf("got invalid response data: %+v", res.Data)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		svc := guestbook.NewService(guestbook.ServiceConfig{})
		ctrl := guestbook.NewController(svc)

		body := []byte(`{invalid-json`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		ctrl.CreateGuestbook(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", w.Code)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		// No context user will trigger service error
		svc := guestbook.NewService(guestbook.ServiceConfig{})
		ctrl := guestbook.NewController(svc)

		body := []byte(`{"text":"Hello world!"}`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		ctrl.CreateGuestbook(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Write Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateGuestbookFn: func(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error) {
				return mockCreateGuestbookRow, nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{
			Database:    db,
			EmailClient: &mockEmailClient{},
		})

		ctrl := guestbook.NewController(svc)

		body := []byte(`{"text":"Hello world!"}`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		ctx := auth.SetUserContext(r.Context(), mockUser)
		r = r.WithContext(ctx)

		w := httptest.NewRecorder()

		errWriter := &test.ErrorResponseRecorder{ResponseRecorder: w}

		ctrl.CreateGuestbook(errWriter, r)

		if w.Code != http.StatusCreated {
			t.Errorf("got %d, want %d", w.Code, http.StatusCreated)
		}
	})
}

func TestController_DeleteGuestbook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			GetGuestbookByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error) {
				return sqlc.Guestbook{ID: id}, nil
			},
			DeleteGuestbookFn: func(ctx context.Context, id pgtype.UUID) error {
				return nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})
		ctrl := guestbook.NewController(svc)

		r := httptest.NewRequest(http.MethodDelete, "/", nil)
		// Set PathValue
		r.SetPathValue("id", uuid.New().String())

		w := httptest.NewRecorder()

		ctrl.DeleteGuestbook(w, r)

		if w.Code != http.StatusNoContent {
			t.Fatalf("got %d, want 204", w.Code)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		svc := guestbook.NewService(guestbook.ServiceConfig{})
		ctrl := guestbook.NewController(svc)

		r := httptest.NewRequest(http.MethodDelete, "/", nil)
		r.SetPathValue("id", "invalid-uuid")

		w := httptest.NewRecorder()

		ctrl.DeleteGuestbook(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", w.Code)
		}
	})
}
