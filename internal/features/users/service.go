package users

import (
	"context"
	"fmt"

	"github.com/ccrsxx/api/internal/db/sqlc"
)

type Service struct {
	db *sqlc.Queries
}

type ServiceConfig struct {
	Database *sqlc.Queries
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db: cfg.Database,
	}
}

func (s *Service) GetListUsers(ctx context.Context) ([]sqlc.User, error) {
	users, err := s.db.ListUsers(ctx)

	if err != nil {
		return nil, fmt.Errorf("list users error: %w", err)
	}

	return users, nil
}
