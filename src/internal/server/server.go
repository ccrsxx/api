package server

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/config"
)

type Server struct {
	port int
}

func NewServer() *http.Server {

	port := config.Env().Port

	server := &Server{
		port: 4000,
	}

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: server.RegisterRoutes(),
	}

	return httpServer
}
