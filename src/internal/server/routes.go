package server

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/features/docs"
	"github.com/ccrsxx/api-go/src/internal/features/home"
	"github.com/ccrsxx/api-go/src/internal/features/tools"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := &api.CustomRouter{ServeMux: http.NewServeMux()}

	home.LoadRoutes(router)
	docs.LoadRoutes(router)
	tools.LoadRoutes(router)

	middlewares := middleware.CreateStack(
		middleware.Cors,
		middleware.Logging,
		middleware.GlobalRateLimit(100, 1*time.Minute),
	)

	return middlewares(router)
}
