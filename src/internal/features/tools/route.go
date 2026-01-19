package tools

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

func LoadRoutes(router *api.CustomRouter) {
	mux := &api.CustomRouter{ServeMux: http.NewServeMux()}

	mux.HandleFunc("GET /ip", GetIpAddress)

	mux.HandleFunc("GET /ipinfo", middleware.HandlerRateLimit(10, 10*time.Second)(GetIpInfo))

	mux.HandleFunc("GET /headers", GetHttpHeaders)

	router.Handle("/tools/", http.StripPrefix("/tools", mux))
}
