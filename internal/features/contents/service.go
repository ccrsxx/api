package contents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/utils"
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

func (s *Service) GetContentsData(ctx context.Context, contentType string) ([]model.Content, error) {
	if contentType != "" {
		if err := utils.Validate.Var(contentType, "content_type"); err != nil {
			return nil, &api.HTTPError{
				Message:    "Invalid content type",
				StatusCode: http.StatusBadRequest,
				Details:    nil,
			}
		}
	}

	dbRows, err := s.db.ListContentByType(ctx, contentType)

	if err != nil {
		return nil, fmt.Errorf("list content by type error: %w", err)
	}

	if dbRows == nil {
		dbRows = []sqlc.ListContentByTypeRow{}
	}

	contents := make([]model.Content, len(dbRows))
	for i, row := range dbRows {
		contents[i] = model.Content{
			Slug:  row.Slug,
			Type:  row.Type,
			Views: row.Views,
			Likes: row.Likes,
		}
	}

	return contents, nil
}

type UpsertContentInput struct {
	Slug string `json:"slug" validate:"required"`
	Type string `json:"type" validate:"required,content_type"`
}

func (s *Service) UpsertContent(ctx context.Context, input UpsertContentInput) (model.Content, error) {
	content, err := s.db.UpsertContent(ctx, sqlc.UpsertContentParams{
		Slug: input.Slug,
		Type: input.Type,
	})

	if err != nil {
		return model.Content{}, fmt.Errorf("upsert content error: %w", err)
	}

	return model.Content{
		Slug:      content.Slug,
		Type:      content.Type,
		CreatedAt: &content.CreatedAt.Time,
		UpdatedAt: &content.UpdatedAt.Time,
	}, nil
}
