package server

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/middleware"
	"github.com/ccrsxx/api-go/src/modules/home"
	"github.com/ccrsxx/api-go/src/modules/tools"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := &api.CustomRouter{ServeMux: http.NewServeMux()}

	home.LoadRoutes(router)
	tools.LoadRoutes(router)

	middlewares := middleware.CreateStack(
		middleware.Cors,
		middleware.Logging,
		middleware.GlobalRateLimit(120, 1*time.Minute),
	)

	return middlewares(router)
}
