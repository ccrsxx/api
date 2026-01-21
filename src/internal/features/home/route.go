package home

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api-go/src/internal/features/tools"
)

func LoadRoutes(router *http.ServeMux) {
	handleHomeRequest := func(w http.ResponseWriter, r *http.Request) {
		hostname := r.Host

		switch {
		case strings.HasPrefix(hostname, "ip."):
			tools.Controller.GetIpAddress(w, r)
			return
		case strings.HasPrefix(hostname, "ipinfo."):
			tools.SharedGetIpInfo.ServeHTTP(w, r)
			return
		case strings.HasPrefix(hostname, "headers."):
			tools.Controller.GetHttpHeaders(w, r)
			return
		}

		Controller.ping(w, r)
	}

	router.HandleFunc("GET /{$}", handleHomeRequest)
}
