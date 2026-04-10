package sse

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/features/auth"
	m "github.com/ccrsxx/api/internal/middleware"
)

type Config struct {
	Router         *http.ServeMux
	Service        *Service
	AuthMiddleware *auth.Middleware
}

func LoadRoutes(config Config) {
	controller := NewController(config.Service)
	middleware := NewMiddleware(config.Service)

	config.Router.Handle("GET /sse",
		config.AuthMiddleware.IsAuthorizedFromQuery(
			m.RateLimit(10, 10*time.Second)(
				middleware.IsConnectionAllowed(
					http.HandlerFunc(controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
