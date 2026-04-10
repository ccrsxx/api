package server

import (
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
)

func NewServer(cfg config.AppConfig) *http.Server {
	InitLoaders(cfg)

	addr := ":" + strconv.Itoa(cfg.Port)

	handler := InitRoutes(cfg)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
