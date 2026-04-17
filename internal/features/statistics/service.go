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
	GetContentStatsByType(ctx context.Context, kind string) (sqlc.GetContentStatsByTypeRow, error)
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

type ContentStatistics struct {
	Type       string `json:"type"`
	TotalPosts int64  `json:"totalPosts"`
	TotalViews int64  `json:"totalViews"`
	TotalLikes int64  `json:"totalLikes"`
}

func (s *Service) GetContentStatistics(ctx context.Context, contentType string) (ContentStatistics, error) {
	if err := utils.Validate.Var(contentType, "content_type"); err != nil {
		return ContentStatistics{}, &api.HTTPError{
			Message:    "Invalid content type",
			StatusCode: http.StatusBadRequest,
		}
	}

	stats, err := s.db.GetContentStatsByType(ctx, contentType)

	if err != nil {
		return ContentStatistics{}, fmt.Errorf("get content stats by type error: %w", err)
	}

	return ContentStatistics{
		Type:       contentType,
		TotalPosts: stats.TotalPosts,
		TotalViews: stats.TotalViews,
		TotalLikes: stats.TotalLikes,
	}, nil
}
