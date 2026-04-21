package guestbook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
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

type Service struct {
	db querier
}

type ServiceConfig struct {
	Database querier
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db: cfg.Database,
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

	return guestbook, nil
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
