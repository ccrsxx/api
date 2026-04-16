package likes

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
	GetContentLikeStatus(ctx context.Context, arg sqlc.GetContentLikeStatusParams) (sqlc.GetContentLikeStatusRow, error)
	CreateContentLike(ctx context.Context, arg sqlc.CreateContentLikeParams) error
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

func (s *Service) getLikeStatus(ctx context.Context, slug string, ipAddressID pgtype.UUID) (sqlc.GetContentLikeStatusRow, error) {
	status, err := s.db.GetContentLikeStatus(ctx, sqlc.GetContentLikeStatusParams{
		Slug:        slug,
		IpAddressID: ipAddressID,
	})

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, fmt.Errorf("get content like status error: %w", err)
	}

	return status, nil
}

func (s *Service) GetLikeStatus(ctx context.Context, slug string, ipAddress string) (sqlc.GetContentLikeStatusRow, error) {
	_, err := s.getContentBySlug(ctx, slug)

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, err
	}

	ip, err := s.db.UpsertIPAddress(ctx, ipAddress)

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, fmt.Errorf("upsert ip address error: %w", err)
	}

	return s.getLikeStatus(ctx, slug, ip.ID)
}

func (s *Service) IncrementLike(ctx context.Context, slug string, ipAddress string) (sqlc.GetContentLikeStatusRow, error) {
	content, err := s.getContentBySlug(ctx, slug)

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, err
	}

	ip, err := s.db.UpsertIPAddress(ctx, ipAddress)

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, fmt.Errorf("upsert ip address error: %w", err)
	}

	status, err := s.getLikeStatus(ctx, slug, ip.ID)

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, err
	}

	if status.UserLikes >= 5 {
		return sqlc.GetContentLikeStatusRow{}, &api.HTTPError{
			Message:    "Likes limit reached",
			StatusCode: http.StatusForbidden,
		}
	}

	err = s.db.CreateContentLike(ctx, sqlc.CreateContentLikeParams{
		ContentID:   content.ID,
		IpAddressID: ip.ID,
	})

	if err != nil {
		return sqlc.GetContentLikeStatusRow{}, fmt.Errorf("create content like error: %w", err)
	}

	return s.getLikeStatus(ctx, slug, ip.ID)
}
