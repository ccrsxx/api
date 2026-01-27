package server

import (
	"net/http"

	"github.com/ccrsxx/api/src/internal/config"
)

func NewServer() *http.Server {
	RegisterLoaders()

	addr := ":" + config.Env().Port

	handler := RegisterRoutes()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
