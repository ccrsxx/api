package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg config.AppConfig, pool *pgxpool.Pool, db *sqlc.Queries) *http.Server {
	LoadLoaders(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)

	handler := LoadHandlers(ctx, cfg, pool, db)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
