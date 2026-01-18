package tools

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

func LoadRoutes(router *api.CustomRouter) {
	r := &api.CustomRouter{ServeMux: http.NewServeMux()}

	r.HandleFunc("GET /ip", getIpAddress)

	r.HandleFunc("GET /ipinfo", middleware.RateLimitFunc(10, 10*time.Second)(getIpInfo))

	r.HandleFunc("GET /headers", getHttpHeaders)

	router.Handle("/tools/", http.StripPrefix("/tools", r))
}
