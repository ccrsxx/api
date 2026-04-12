package sse

import (
	"context"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/middleware"
)

type Config struct {
	Router         *http.ServeMux
	Service        *Service
	AppContext     context.Context
	AuthMiddleware *auth.Middleware
}

func LoadRoutes(cfg Config) {
	mw := NewMiddleware(cfg.Service)
	ctrl := NewController(cfg.AppContext, cfg.Service)

	cfg.Router.Handle("GET /sse",
		cfg.AuthMiddleware.IsAuthorizedFromQuery(
			middleware.RateLimit(cfg.AppContext, 10, 10*time.Second)(
				mw.IsConnectionAllowed(
					http.HandlerFunc(ctrl.getCurrentPlayingSSE),
				),
			),
		),
	)
}
