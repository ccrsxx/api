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

func LoadRoutes(config Config) {
	controller := NewController()

	handleHomeRequest := func(w http.ResponseWriter, r *http.Request) {
		hostname := r.Host

		switch {
		case strings.HasPrefix(hostname, "ip."):
			config.ToolsController.GetIpAddress(w, r)
			return
		case strings.HasPrefix(hostname, "ipinfo."):
			config.SharedGetIpInfoController.ServeHTTP(w, r)
			return
		case strings.HasPrefix(hostname, "headers."):
			config.ToolsController.GetHttpHeaders(w, r)
			return
		}

		controller.ping(w, r)
	}

	config.Router.HandleFunc("GET /{$}", handleHomeRequest)
}
