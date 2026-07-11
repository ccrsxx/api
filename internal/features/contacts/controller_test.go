package contacts_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/contacts"
	"github.com/ccrsxx/api/internal/test"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestController_CreateContact(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateContactFn: func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
				return mockContact, nil
			},
			UpdateContactDeliveredAtByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Contact, error) {
				return mockContact, nil
			},
		}

		svc := contacts.NewService(contacts.ServiceConfig{
			Database:         db,
			PushoverClient:   &mockPushoverClient{},
			EmailClient:      &mockEmailClient{},
			CloudflareClient: &mockCloudflareClient{},
		})

		ctrl := contacts.NewController(svc)

		body := []byte(`{"name":"Test","email":"test@example.com","message":"Hello","token":"test_token"}`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		ctrl.CreateContact(w, r)

		if w.Code != http.StatusNoContent {
			t.Fatalf("got %d, want 204", w.Code)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		svc := contacts.NewService(contacts.ServiceConfig{})
		ctrl := contacts.NewController(svc)

		body := []byte(`{invalid-json`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		ctrl.CreateContact(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", w.Code)
		}
	})

	t.Run("Turnstile Error", func(t *testing.T) {
		svc := contacts.NewService(contacts.ServiceConfig{
			CloudflareClient: &mockCloudflareClient{err: errors.New("turnstile error")},
		})

		ctrl := contacts.NewController(svc)

		body := []byte(`{"name":"Test","email":"test@example.com","message":"Hello","token":"test_token"}`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		ctrl.CreateContact(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateContactFn: func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
				return sqlc.Contact{}, errors.New("db error")
			},
		}

		svc := contacts.NewService(contacts.ServiceConfig{
			Database:         db,
			CloudflareClient: &mockCloudflareClient{},
		})

		ctrl := contacts.NewController(svc)

		body := []byte(`{"name":"Test","email":"test@example.com","message":"Hello","token":"test_token"}`)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		ctrl.CreateContact(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want 500", w.Code)
		}
	})
}
