package guestbook

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/clients/gmail"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/google/uuid"
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

type emailClient interface {
	Send(msg gmail.Message) error
}

type Service struct {
	db           querier
	emailClient  emailClient
	emailTarget  string
	emailAddress string
}

type ServiceConfig struct {
	Database     querier
	EmailClient  emailClient
	EmailTarget  string
	EmailAddress string
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db:           cfg.Database,
		emailClient:  cfg.EmailClient,
		emailTarget:  cfg.EmailTarget,
		emailAddress: cfg.EmailAddress,
	}
}

type CreateGuestbookInput struct {
	Text string `json:"text" validate:"required"`
}

func (s *Service) CreateGuestbook(ctx context.Context, input CreateGuestbookInput) (sqlc.CreateGuestbookRow, error) {
	user, err := auth.GetUserFromContext(ctx)

	if err != nil {
		return sqlc.CreateGuestbookRow{}, fmt.Errorf("get user by id error: %w", err)
	}

	guestbook, err := s.db.CreateGuestbook(ctx, sqlc.CreateGuestbookParams{
		Text:   input.Text,
		UserID: user.ID,
	})

	if err != nil {
		return sqlc.CreateGuestbookRow{}, fmt.Errorf("create guestbook error: %w", err)
	}

	go s.sendNewGuestbookEmail(user, guestbook)

	return guestbook, nil
}

func (s *Service) sendNewGuestbookEmail(user sqlc.GetUserWithAccountByIDRow, guestbook sqlc.CreateGuestbookRow) {
	subject := fmt.Sprintf("New guestbook from %s (%s)", user.Name, user.Email.String)

	err := s.emailClient.Send(gmail.Message{
		From:    s.emailAddress,
		To:      s.emailTarget,
		Subject: subject,
		Text:    guestbook.Text,
	})

	if err != nil {
		slog.Warn("send new guestbook email error", "error", err)
	}
}

func (s *Service) ListGuestbook(ctx context.Context) ([]sqlc.ListGuestbookRow, error) {
	guestbook, err := s.db.ListGuestbook(ctx)

	if err != nil {
		return nil, fmt.Errorf("list guestbook error: %w", err)
	}

	if guestbook == nil {
		guestbook = []sqlc.ListGuestbookRow{}
	}

	return guestbook, nil
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
