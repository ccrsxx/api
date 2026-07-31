package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New builds the public API server.
//
// LoadLoaders must have been called by the caller beforehand, so that logging
// and tracing are configured before any instrumented component (the pgx pool,
// the HTTP handlers) captures the global tracer provider.
func New(ctx context.Context, cfg config.AppConfig, pool *pgxpool.Pool, db *sqlc.Queries) *http.Server {
	addr := ":" + strconv.Itoa(cfg.Port)

	handler := LoadHandlers(ctx, cfg, pool, db)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
