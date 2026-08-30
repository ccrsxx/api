package guestbook

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"uuid"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/gmail"
	"github.com/ccrsxx/api/internal/clients/pushover"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type querier interface {
	GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error)
	GetGuestbookByID(ctx context.Context, id pgtype.UUID) (sqlc.Guestbook, error)
	ListGuestbook(ctx context.Context) ([]sqlc.ListGuestbookRow, error)
	CreateGuestbook(ctx context.Context, arg sqlc.CreateGuestbookParams) (sqlc.CreateGuestbookRow, error)
	DeleteGuestbook(ctx context.Context, id pgtype.UUID) error
}

type pushoverClient interface {
	SendMessage(ctx context.Context, messageRequest pushover.MessageRequest) error
}

type emailClient interface {
	Send(msg gmail.Message) error
}

type Service struct {
	db             querier
	emailClient    emailClient
	emailTarget    string
	emailAddress   string
	pushoverClient pushoverClient
}

type ServiceConfig struct {
	Database       querier
	EmailClient    emailClient
	EmailTarget    string
	EmailAddress   string
	PushoverClient pushoverClient
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db:             cfg.Database,
		emailClient:    cfg.EmailClient,
		emailTarget:    cfg.EmailTarget,
		emailAddress:   cfg.EmailAddress,
		pushoverClient: cfg.PushoverClient,
	}
}

type CreateGuestbookInput struct {
	Text string `json:"text" validate:"required"`
}

func (s *Service) CreateGuestbook(ctx context.Context, input CreateGuestbookInput) (Guestbook, error) {
	user, err := auth.GetUserFromContext(ctx)

	if err != nil {
		return Guestbook{}, fmt.Errorf("get user by id error: %w", err)
	}

	guestbook, err := s.db.CreateGuestbook(ctx, sqlc.CreateGuestbookParams{
		Text:   input.Text,
		UserID: user.ID,
	})

	if err != nil {
		return Guestbook{}, fmt.Errorf("create guestbook error: %w", err)
	}

	go s.sendNewGuestbookNotifications(user, guestbook)

	return Guestbook{
		ID:        guestbook.ID.String(),
		Text:      guestbook.Text,
		Name:      guestbook.Name,
		Image:     guestbook.Image.String,
		Username:  guestbook.Username.String,
		CreatedAt: guestbook.CreatedAt.Time,
	}, nil
}

func (s *Service) sendNewGuestbookNotifications(user sqlc.GetUserWithAccountByIDRow, guestbook sqlc.CreateGuestbookRow) {
	subject := fmt.Sprintf("New guestbook from %s (%s)", user.Name, user.Email.String)

	err := s.pushoverClient.SendMessage(context.Background(), pushover.MessageRequest{
		Title:   subject,
		Message: guestbook.Text,
	})

	// Log the error but don't return it, since this is a background task
	if err != nil {
		slog.Error("send new guestbook pushover error", "error", err)
	}

	err = s.emailClient.Send(gmail.Message{
		From:    s.emailAddress,
		To:      s.emailTarget,
		Subject: subject,
		Text:    guestbook.Text,
	})

	// Log the error but don't return it, since this is a background task
	if err != nil {
		slog.Error("send new guestbook email error", "error", err)
	}
}

func (s *Service) ListGuestbook(ctx context.Context) ([]Guestbook, error) {
	dbRows, err := s.db.ListGuestbook(ctx)

	if err != nil {
		return nil, fmt.Errorf("list guestbook error: %w", err)
	}

	guestbooks := make([]Guestbook, len(dbRows))

	for i, row := range dbRows {
		guestbooks[i] = Guestbook{
			ID:        row.ID.String(),
			Text:      row.Text,
			Name:      row.Name,
			Image:     row.Image.String,
			Username:  row.Username.String,
			CreatedAt: row.CreatedAt.Time,
		}
	}

	return guestbooks, nil
}

func (s *Service) DeleteGuestbook(ctx context.Context, guestbookID string) error {
	id, err := uuid.Parse(guestbookID)

	if err != nil {
		return &api.HTTPError{
			Message:    "Invalid id",
			StatusCode: http.StatusBadRequest,
		}
	}

	guestbook, err := s.db.GetGuestbookByID(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return &api.HTTPError{
			Message:    "Guestbook not found",
			StatusCode: http.StatusNotFound,
		}
	}

	if err != nil {
		return fmt.Errorf("get guestbook by id error: %w", err)
	}

	if err := s.db.DeleteGuestbook(ctx, guestbook.ID); err != nil {
		return fmt.Errorf("delete guestbook error: %w", err)
	}

	return nil
}
