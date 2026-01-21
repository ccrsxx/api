package tools

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/middleware"
)

// Shared rate-limited handler for GetIpInfo. Limits to 10 requests per 10 seconds.
var SharedGetIpInfo = middleware.HandlerRateLimit(10, 10*time.Second)(
	http.HandlerFunc(Controller.GetIpInfo),
)

func LoadRoutes(router *http.ServeMux) {
	mux := http.NewServeMux()

	mux.Handle("GET /ipinfo", SharedGetIpInfo)

	mux.HandleFunc("GET /ip", Controller.GetIpInfo)

	mux.HandleFunc("GET /headers", Controller.GetHttpHeaders)

	router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
