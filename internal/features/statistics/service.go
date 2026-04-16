package statistics

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/db/sqlc"
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

type ContentStatistics struct {
	Type       string `json:"type"`
	TotalPosts int32  `json:"totalPosts"`
	TotalViews int32  `json:"totalViews"`
	TotalLikes int32  `json:"totalLikes"`
}

func (s *Service) GetContentStatistics(ctx context.Context, contentType string) (ContentStatistics, error) {
	if err := validateContentType(contentType); err != nil {
		return ContentStatistics{}, err
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
