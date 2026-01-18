package server

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/modules/home"
	"github.com/ccrsxx/api-go/src/modules/tools"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := &api.CustomRouter{ServeMux: http.NewServeMux()}

	home.LoadRoutes(router)
	tools.LoadRoutes(router)

	return router
}
