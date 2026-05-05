package guestbook_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/gmail"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/guestbook"
	"github.com/ccrsxx/api/internal/test"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockEmailClient struct {
	err error
}

func (m *mockEmailClient) Send(msg gmail.Message) error {
	return m.err
}

var (
	mockUserID      = pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	mockGuestbookID = pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockUser        = sqlc.GetUserWithAccountByIDRow{
		ID:   mockUserID,
		Name: "Test User",
		Email: pgtype.Text{
			String: "test@example.com",
			Valid:  true,
		},
	}
	mockCreateGuestbookRow = sqlc.CreateGuestbookRow{
		ID:   mockGuestbookID,
		Text: "Hello world!",
		Name: "Test User",
	}
	mockListGuestbookRow = []sqlc.ListGuestbookRow{
		{
			ID:   mockGuestbookID,
			Text: "Hello world!",
			Name: "Test User",
		},
	}
)

func TestService_CreateGuestbook(t *testing.T) {
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

		ctx := auth.SetUserContext(context.Background(), mockUser)

		input := guestbook.CreateGuestbookInput{Text: "Hello world!"}

		result, err := svc.CreateGuestbook(ctx, input)

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if result.ID != mockCreateGuestbookRow.ID {
			t.Errorf("got %v, want %v", result.ID, mockCreateGuestbookRow.ID)
		}
	})

	t.Run("Context User Error", func(t *testing.T) {
		svc := guestbook.NewService(guestbook.ServiceConfig{})

		_, err := svc.CreateGuestbook(context.Background(), guestbook.CreateGuestbookInput{Text: "Hello"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get user by id error") {
			t.Errorf("got %v, want get user by id error", err)
		}
	})

	t.Run("Create Guestbook DB Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateGuestbookFn: func(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error) {
				return sqlc.CreateGuestbookRow{}, errors.New("db error")
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{
			Database: db,
		})

		ctx := auth.SetUserContext(context.Background(), mockUser)

		_, err := svc.CreateGuestbook(ctx, guestbook.CreateGuestbookInput{Text: "Hello"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "create guestbook error") {
			t.Errorf("got %v, want create guestbook error", err)
		}
	})

	t.Run("Email Error", func(t *testing.T) {
		db := &test.MockQuerier{
			CreateGuestbookFn: func(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error) {
				return mockCreateGuestbookRow, nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{
			Database:    db,
			EmailClient: &mockEmailClient{err: errors.New("email fail")},
		})

		ctx := auth.SetUserContext(context.Background(), mockUser)

		_, err := svc.CreateGuestbook(ctx, guestbook.CreateGuestbookInput{Text: "Hello"})

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		// Wait for goroutine to hit the slog.Error path.
		time.Sleep(10 * time.Millisecond)
	})
}

func TestService_ListGuestbook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := &test.MockQuerier{
			ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
				return mockListGuestbookRow, nil
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

		list, err := svc.ListGuestbook(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if len(list) != 1 {
			t.Fatalf("got %d, want 1", len(list))
		}
	})

	t.Run("Fallback to Empty Slice", func(t *testing.T) {
		db := &test.MockQuerier{
			ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
				return nil, nil // Return nil slice to test fallback
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

		list, err := svc.ListGuestbook(context.Background())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if list == nil {
			t.Fatal("expected empty slice, got nil")
		}
	})

	t.Run("DB Error", func(t *testing.T) {
		db := &test.MockQuerier{
			ListGuestbookFn: func(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
				return nil, errors.New("db error")
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

		_, err := svc.ListGuestbook(context.Background())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "list guestbook error") {
			t.Errorf("got %v, want list guestbook error", err)
		}
	})
}

func TestService_DeleteGuestbook(t *testing.T) {
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

		err := svc.DeleteGuestbook(context.Background(), uuid.New().String())

		if err != nil {
			t.Fatalf("unwanted error: %v", err)
		}
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := guestbook.NewService(guestbook.ServiceConfig{})

		err := svc.DeleteGuestbook(context.Background(), "invalid-uuid")

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want *api.HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %v", err)
		}
	})

	t.Run("Guestbook Not Found", func(t *testing.T) {
		db := &test.MockQuerier{
			GetGuestbookByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error) {
				return sqlc.Guestbook{}, pgx.ErrNoRows
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

		err := svc.DeleteGuestbook(context.Background(), uuid.New().String())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		httpErr, ok := errors.AsType[*api.HTTPError](err)

		if !ok {
			t.Fatalf("got %T, want *api.HTTPError", err)
		}

		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 Not Found, got %v", err)
		}
	})

	t.Run("Get Guestbook Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetGuestbookByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error) {
				return sqlc.Guestbook{}, errors.New("db error")
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

		err := svc.DeleteGuestbook(context.Background(), uuid.New().String())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "get guestbook by id error") {
			t.Errorf("got %v, want get guestbook by id error", err)
		}
	})

	t.Run("Delete DB Error", func(t *testing.T) {
		db := &test.MockQuerier{
			GetGuestbookByIDFn: func(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error) {
				return sqlc.Guestbook{ID: id}, nil
			},
			DeleteGuestbookFn: func(ctx context.Context, id pgtype.UUID) error {
				return errors.New("db delete error")
			},
		}

		svc := guestbook.NewService(guestbook.ServiceConfig{Database: db})

		err := svc.DeleteGuestbook(context.Background(), uuid.New().String())

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "delete guestbook error") {
			t.Errorf("got %v, want delete guestbook error", err)
		}
	})
}
