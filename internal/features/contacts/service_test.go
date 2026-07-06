package contacts_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/clients/gmail"
	"github.com/ccrsxx/api/internal/clients/pushover"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/contacts"
	"github.com/ccrsxx/api/internal/test"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockCloudflareClient struct {
	err error
}

func (m *mockCloudflareClient) VerifyTurnstile(ctx context.Context, token string, remoteIP string) error {
	return m.err
}

type mockPushoverClient struct {
	err error
}

func (m *mockPushoverClient) SendMessage(ctx context.Context, messageRequest pushover.MessageRequest) error {
	return m.err
}

type mockEmailClient struct {
	err error
}

func (m *mockEmailClient) Send(msg gmail.Message) error {
	return m.err
}

var (
	mockContactID = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockContact   = sqlc.Contact{
		ID:      mockContactID,
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Hello!",
	}
)

func TestService_CreateContact(t *testing.T) {
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
			Database:       db,
			PushoverClient: &mockPushoverClient{},
			EmailClient:    &mockEmailClient{},
		})

		input := contacts.CreateContactInput{
			Name:    "Test User",
			Email:   "test@example.com",
			Message: "Hello!",
			Token:   "test_token",
		}

		err := svc.CreateContact(context.Background(), input, "127.0.0.1")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}
	})

	t.Run("Create Contact DB Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateContactFn: func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
				return sqlc.Contact{}, errors.New("db error")
			},
		}

		svc := contacts.NewService(contacts.ServiceConfig{Database: db})

		err := svc.CreateContact(context.Background(), contacts.CreateContactInput{}, "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create contact create record error") {
			t.Errorf("got %v, want create contact create record error", err)
		}
	})

	t.Run("Pushover Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateContactFn: func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
				return mockContact, nil
			},
		}

		svc := contacts.NewService(contacts.ServiceConfig{
			Database:       db,
			PushoverClient: &mockPushoverClient{err: errors.New("pushover error")},
		})

		err := svc.CreateContact(context.Background(), contacts.CreateContactInput{}, "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create contact pushover notification error") {
			t.Errorf("got %v, want create contact pushover notification error", err)
		}
	})

	t.Run("Update Delivered DB Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateContactFn: func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
				return mockContact, nil
			},
			UpdateContactDeliveredAtByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Contact, error) {
				return sqlc.Contact{}, errors.New("db update error")
			},
		}

		svc := contacts.NewService(contacts.ServiceConfig{
			Database:       db,
			PushoverClient: &mockPushoverClient{},
			EmailClient:    &mockEmailClient{},
		})

		err := svc.CreateContact(context.Background(), contacts.CreateContactInput{}, "127.0.0.1")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create contact update record error") {
			t.Errorf("got %v, want create contact update record error", err)
		}
	})

	t.Run("Email Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateContactFn: func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
				return mockContact, nil
			},
			UpdateContactDeliveredAtByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Contact, error) {
				return mockContact, nil
			},
		}

		svc := contacts.NewService(contacts.ServiceConfig{
			Database:       db,
			PushoverClient: &mockPushoverClient{},
			EmailClient:    &mockEmailClient{err: errors.New("email fail")},
		})

		err := svc.CreateContact(context.Background(), contacts.CreateContactInput{}, "127.0.0.1")

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		// Wait for goroutine to hit the slog.Error path.
		time.Sleep(10 * time.Millisecond)
	})
}
