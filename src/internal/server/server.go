package server

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	config.LoadEnv()

	port := config.Env().Port

	server := &Server{
		port: 4000,
	}

	middlewares := middleware.CreateStack(
		middleware.Logging,
	)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: middlewares(server.RegisterRoutes()),
	}

	return httpServer
}
