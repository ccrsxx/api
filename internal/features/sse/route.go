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

func LoadRoutes(cfg Config) {
	controller := NewController(cfg.Service)
	middleware := NewMiddleware(cfg.Service)

	cfg.Router.Handle("GET /sse",
		cfg.AuthMiddleware.IsAuthorizedFromQuery(
			m.RateLimit(10, 10*time.Second)(
				middleware.IsConnectionAllowed(
					http.HandlerFunc(controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
