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
	ctrl := NewController(cfg.Service)
	mw := NewMiddleware(cfg.Service)

	cfg.Router.Handle("GET /sse",
		cfg.AuthMiddleware.IsAuthorizedFromQuery(
			m.RateLimit(10, 10*time.Second)(
				mw.IsConnectionAllowed(
					http.HandlerFunc(ctrl.getCurrentPlayingSSE),
				),
			),
		),
	)
}
