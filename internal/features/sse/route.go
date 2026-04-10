package sse

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/features/auth"
	m "github.com/ccrsxx/api/internal/middleware"
)

func LoadRoutes(router *http.ServeMux, service *Service, authMiddleware *auth.Middleware) {
	controller := NewController(service)
	middleware := NewMiddleware(service)

	router.Handle("GET /sse",
		authMiddleware.IsAuthorizedFromQuery(
			m.RateLimit(10, 10*time.Second)(
				middleware.IsConnectionAllowed(
					http.HandlerFunc(controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
