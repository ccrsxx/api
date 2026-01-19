package tools

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

// Shared rate-limited handler for GetIpInfo. Limits to 10 requests per 10 seconds.
var SharedGetIpInfo = middleware.HandlerRateLimit(1, 10*time.Second)(GetIpInfo)

func LoadRoutes(router *api.CustomRouter) {
	mux := &api.CustomRouter{ServeMux: http.NewServeMux()}

	mux.HandleFunc("GET /ip", GetIpAddress)

	mux.HandleFunc("GET /ipinfo", SharedGetIpInfo)

	mux.HandleFunc("GET /headers", GetHttpHeaders)

	router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
