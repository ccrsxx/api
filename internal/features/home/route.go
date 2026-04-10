package home

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api/internal/features/tools"
)

func LoadRoutes(router *http.ServeMux, config tools.Config) {
	controller := NewController()

	handleHomeRequest := func(w http.ResponseWriter, r *http.Request) {
		hostname := r.Host

		switch {
		case strings.HasPrefix(hostname, "ip."):
			config.ToolsController.GetIpAddress(w, r)
			return
		case strings.HasPrefix(hostname, "ipinfo."):
			config.SharedGetIpInfo.ServeHTTP(w, r)
			return
		case strings.HasPrefix(hostname, "headers."):
			config.ToolsController.GetHttpHeaders(w, r)
			return
		}

		controller.ping(w, r)
	}

	router.HandleFunc("GET /{$}", handleHomeRequest)
}
