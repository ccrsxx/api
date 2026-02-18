package server

import (
	"net/http"
	"strconv"

	"github.com/ccrsxx/api/internal/config"
)

func NewServer() *http.Server {
	RegisterLoaders()

	addr := ":" + strconv.Itoa(config.Env().Port)

	handler := RegisterRoutes()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return httpServer
}
