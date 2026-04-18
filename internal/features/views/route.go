package views

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
	mux := http.NewServeMux()

	ctrl := NewController(cfg.Service)

	mux.HandleFunc("GET /{slug}", ctrl.GetViewCount)

	mux.Handle("POST /{slug}",
		cfg.AuthMiddleware.IsAuthorized(
			middleware.RateLimit(cfg.AppContext, 60, 1*time.Hour)(
				http.HandlerFunc(ctrl.IncrementView),
			),
		),
	)

	cfg.Router.Handle("/views/", http.StripPrefix("/views", mux))
}
