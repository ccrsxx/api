package statistics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/utils"
)

type querier interface {
	GetContentStatsByType(ctx context.Context, type_ string) (sqlc.GetContentStatsByTypeRow, error)
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

func (s *Service) GetContentsStatistics(ctx context.Context, contentType string) (sqlc.GetContentStatsByTypeRow, error) {
	if contentType != "" {
		if err := utils.Validate.Var(contentType, "content_type"); err != nil {
			return sqlc.GetContentStatsByTypeRow{}, &api.HTTPError{
				Message:    "Invalid content type",
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	stats, err := s.db.GetContentStatsByType(ctx, contentType)

	if err != nil {
		return sqlc.GetContentStatsByTypeRow{}, fmt.Errorf("get content stats by type error: %w", err)
	}

	return sqlc.GetContentStatsByTypeRow{
		TotalPosts: stats.TotalPosts,
		TotalViews: stats.TotalViews,
		TotalLikes: stats.TotalLikes,
	}, nil
}
