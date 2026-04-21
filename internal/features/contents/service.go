package contents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
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

func (s *Service) GetContentsData(ctx context.Context, contentType string) ([]sqlc.ListContentByTypeRow, error) {
	if contentType != "" {
		if err := utils.Validate.Var(contentType, "content_type"); err != nil {
			return nil, &api.HTTPError{
				Message:    "Invalid content type",
				StatusCode: http.StatusBadRequest,
				Details:    nil,
			}
		}
	}

	data, err := s.db.ListContentByType(ctx, contentType)

	if err != nil {
		return nil, fmt.Errorf("list content by type error: %w", err)
	}

	if data == nil {
		data = []sqlc.ListContentByTypeRow{}
	}

	return data, nil
}

type UpsertContentInput struct {
	Slug string `json:"slug" validate:"required"`
	Type string `json:"type" validate:"required,content_type"`
}

func (s *Service) UpsertContent(ctx context.Context, input UpsertContentInput) (sqlc.Content, error) {
	content, err := s.db.UpsertContent(ctx, sqlc.UpsertContentParams{
		Slug: input.Slug,
		Type: input.Type,
	})

	if err != nil {
		return sqlc.Content{}, fmt.Errorf("upsert content error: %w", err)
	}

	return content, nil
}
