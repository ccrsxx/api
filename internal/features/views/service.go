package views

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type querier interface {
	GetContentBySlug(ctx context.Context, slug string) (sqlc.Content, error)
	GetTotalContentMeta(ctx context.Context, contentID pgtype.UUID) (sqlc.GetTotalContentMetaRow, error)
	IncrementContentView(ctx context.Context, arg sqlc.IncrementContentViewParams) (sqlc.IncrementContentViewRow, error)
	UpsertIPAddress(ctx context.Context, ipAddress string) (sqlc.IpAddress, error)
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

type ViewCount struct {
	Views int64 `json:"views"`
}

func (s *Service) getContentBySlug(ctx context.Context, slug string) (sqlc.Content, error) {
	content, err := s.db.GetContentBySlug(ctx, slug)

	if errors.Is(err, pgx.ErrNoRows) {
		return sqlc.Content{}, &api.HTTPError{
			Message:    "Content not found",
			StatusCode: http.StatusNotFound,
		}
	}

	if err != nil {
		return sqlc.Content{}, fmt.Errorf("get content by slug error: %w", err)
	}

	return content, nil
}

func (s *Service) GetViewCount(ctx context.Context, slug string) (ViewCount, error) {
	content, err := s.getContentBySlug(ctx, slug)

	if err != nil {
		return ViewCount{}, err
	}

	meta, err := s.db.GetTotalContentMeta(ctx, content.ID)

	if err != nil {
		return ViewCount{}, fmt.Errorf("get view count error: %w", err)
	}

	return ViewCount{Views: meta.TotalViews}, nil
}

func (s *Service) IncrementView(ctx context.Context, slug string, ipAddress string) (ViewCount, error) {
	/*
		TODO: Better approach maybe be in the future with less call to db, it's possible to apply to /likes too
		Get content with full stats -> Increment -> Increment from the first struct
	*/

	content, err := s.getContentBySlug(ctx, slug)

	if err != nil {
		return ViewCount{}, err
	}

	ip, err := s.db.UpsertIPAddress(ctx, ipAddress)

	if err != nil {
		return ViewCount{}, fmt.Errorf("upsert ip address error: %w", err)
	}

	_, err = s.db.IncrementContentView(ctx, sqlc.IncrementContentViewParams{
		ContentID:   content.ID,
		IpAddressID: ip.ID,
	})

	if err != nil {
		return ViewCount{}, fmt.Errorf("create content view error: %w", err)
	}

	meta, err := s.db.GetTotalContentMeta(ctx, content.ID)

	if err != nil {
		return ViewCount{}, fmt.Errorf("get view count after increment error: %w", err)
	}

	return ViewCount{Views: meta.TotalViews}, nil
}
