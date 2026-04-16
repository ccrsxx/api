package views

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
)

type querier interface {
	GetContentBySlug(ctx context.Context, slug string) (sqlc.Content, error)
	GetContentViewCount(ctx context.Context, slug string) (int32, error)
	CreateContentView(ctx context.Context, arg sqlc.CreateContentViewParams) error
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
	Views int32 `json:"views"`
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
	_, err := s.getContentBySlug(ctx, slug)

	if err != nil {
		return ViewCount{}, err
	}

	views, err := s.db.GetContentViewCount(ctx, slug)

	if err != nil {
		return ViewCount{}, fmt.Errorf("get view count error: %w", err)
	}

	return ViewCount{Views: views}, nil
}

func (s *Service) IncrementView(ctx context.Context, slug string, ipAddress string) (ViewCount, error) {
	content, err := s.getContentBySlug(ctx, slug)

	if err != nil {
		return ViewCount{}, err
	}

	ip, err := s.db.UpsertIPAddress(ctx, ipAddress)

	if err != nil {
		return ViewCount{}, fmt.Errorf("upsert ip address error: %w", err)
	}

	err = s.db.CreateContentView(ctx, sqlc.CreateContentViewParams{
		ContentID:   content.ID,
		IpAddressID: ip.ID,
	})

	if err != nil {
		return ViewCount{}, fmt.Errorf("create content view error: %w", err)
	}

	views, err := s.db.GetContentViewCount(ctx, slug)

	if err != nil {
		return ViewCount{}, fmt.Errorf("get view count after increment error: %w", err)
	}

	return ViewCount{Views: views}, nil
}
