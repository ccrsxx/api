package tools

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

func LoadRoutes(r *api.CustomRouter) {
	mux := &api.CustomRouter{ServeMux: http.NewServeMux()}

	mux.HandleFunc("GET /ip", getIpAddress)

	mux.HandleFunc("GET /ipinfo", middleware.RateLimitFunc(10, 10*time.Second)(getIpInfo))

	mux.HandleFunc("GET /headers", getHttpHeaders)

	r.Handle("/tools/", http.StripPrefix("/tools", mux))
}
