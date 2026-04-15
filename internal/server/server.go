package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
)

func New(ctx context.Context, cfg config.AppConfig, db *sqlc.Queries) *http.Server {
	LoadLoaders(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)

	handler := LoadHandlers(ctx, cfg, db)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
