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

type ContentsStatistics struct {
	Type       string `json:"type"`
	TotalPosts int64  `json:"totalPosts"`
	TotalViews int64  `json:"totalViews"`
	TotalLikes int64  `json:"totalLikes"`
}

func (s *Service) GetContentsStatistics(ctx context.Context, contentType string) (ContentsStatistics, error) {
	if contentType != "" {
		if err := utils.Validate.Var(contentType, "content_type"); err != nil {
			return ContentsStatistics{}, &api.HTTPError{
				Message:    "Invalid content type",
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	stats, err := s.db.GetContentStatsByType(ctx, contentType)

	if err != nil {
		return ContentsStatistics{}, fmt.Errorf("get content stats by type error: %w", err)
	}

	parsedContentType := contentType

	if parsedContentType == "" {
		parsedContentType = "all"
	}

	return ContentsStatistics{
		Type:       parsedContentType,
		TotalPosts: stats.TotalPosts,
		TotalViews: stats.TotalViews,
		TotalLikes: stats.TotalLikes,
	}, nil
}
