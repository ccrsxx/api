package contacts

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
	ctrl := NewController(cfg.Service)

	cfg.Router.Handle("POST /contacts",
		cfg.AuthMiddleware.IsAuthorizedFromBearer(
			middleware.RateLimit(cfg.AppContext, 20, 1*time.Hour)(
				http.HandlerFunc(ctrl.CreateContact),
			),
		),
	)
}
