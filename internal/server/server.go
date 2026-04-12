package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
)

func New(ctx context.Context, cfg config.AppConfig) *http.Server {
	LoadLoaders(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)

	handler := LoadHandlers(ctx, cfg)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
