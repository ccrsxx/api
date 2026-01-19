package home

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/features/tools" // Import your tools package
)

func LoadRoutes(router *api.CustomRouter) {
	handleHomeRequest := func(w http.ResponseWriter, r *http.Request) error {
		hostname := r.Host

		switch {
		case strings.HasPrefix(hostname, "ip."):
			return tools.GetIpAddress(w, r)
		case strings.HasPrefix(hostname, "ipinfo."):
			return tools.GetIpInfo(w, r)
		case strings.HasPrefix(hostname, "headers."):
			return tools.GetHttpHeaders(w, r)
		}

		return ping(w, r)
	}

	router.HandleFunc("GET /{$}", handleHomeRequest)
}
