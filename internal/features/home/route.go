package home

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api/internal/features/tools"
)

type Config struct {
	Router                    *http.ServeMux
	ToolsController           *tools.Controller
	SharedGetIpInfoController http.Handler
}

func LoadRoutes(cfg Config) {
	ctrl := NewController()

	handleHomeRequest := func(w http.ResponseWriter, r *http.Request) {
		hostname := r.Host

		switch {
		case strings.HasPrefix(hostname, "ip."):
			cfg.ToolsController.GetIpAddress(w, r)
			return
		case strings.HasPrefix(hostname, "ipinfo."):
			cfg.SharedGetIpInfoController.ServeHTTP(w, r)
			return
		case strings.HasPrefix(hostname, "headers."):
			cfg.ToolsController.GetHttpHeaders(w, r)
			return
		}

		ctrl.ping(w, r)
	}

	cfg.Router.HandleFunc("GET /{$}", handleHomeRequest)
}
