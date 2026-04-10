package server

import (
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
)

func NewServer(cfg config.AppConfig) *http.Server {
	RegisterLoaders(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)

	handler := RegisterRoutes(cfg)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
