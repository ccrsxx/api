package server

import (
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
)

func New(cfg config.AppConfig) *http.Server {
	LoadLoaders(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)

	handler := LoadHandlers(cfg)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
