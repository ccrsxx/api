package sse

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/features/auth"
	m "github.com/ccrsxx/api/internal/middleware"
)

func LoadRoutes(router *http.ServeMux, svc *service) {
	controller := NewController(svc)
	middleware := NewMiddleware(svc)

	router.Handle("GET /sse",
		auth.Middleware.IsAuthorizedFromQuery(
			m.RateLimit(10, 10*time.Second)(
				middleware.IsConnectionAllowed(
					http.HandlerFunc(controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
