package tools

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/api"
)

func LoadRoutes(router *api.CustomRouter) {
	r := &api.CustomRouter{ServeMux: http.NewServeMux()}

	r.HandleFunc("GET /ip", getIpAddress)
	r.HandleFunc("GET /ipinfo", getIpInfo)
	r.HandleFunc("GET /headers", getHttpHeaders)

	router.Handle("/tools/", http.StripPrefix("/tools", r))
}
