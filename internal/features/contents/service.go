package contents

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
)

type querier interface {
	ListContentByType(ctx context.Context, type_ string) ([]sqlc.ListContentByTypeRow, error)
	UpsertContent(ctx context.Context, arg sqlc.UpsertContentParams) (sqlc.Content, error)
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

var validContentTypes = []string{"blog", "project"}

func validateContentType(contentType string) error {
	if !slices.Contains(validContentTypes, contentType) {
		return &api.HTTPError{
			Message:    "Invalid content type",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}

type ContentData struct {
	Type string                      `json:"type"`
	Data []sqlc.ListContentByTypeRow `json:"data"`
}

func (s *Service) GetContentData(ctx context.Context, contentType string) (ContentData, error) {
	if err := validateContentType(contentType); err != nil {
		return ContentData{}, err
	}

	data, err := s.db.ListContentByType(ctx, contentType)

	if err != nil {
		return ContentData{}, fmt.Errorf("list content by type error: %w", err)
	}

	if data == nil {
		data = []sqlc.ListContentByTypeRow{}
	}

	return ContentData{
		Type: contentType,
		Data: data,
	}, nil
}

func (s *Service) UpsertContent(ctx context.Context, slug string, contentType string) (sqlc.Content, error) {
	if err := validateContentType(contentType); err != nil {
		return sqlc.Content{}, err
	}

	content, err := s.db.UpsertContent(ctx, sqlc.UpsertContentParams{
		Slug: slug,
		Type: contentType,
	})

	if err != nil {
		return sqlc.Content{}, fmt.Errorf("upsert content error: %w", err)
	}

	return content, nil
}
