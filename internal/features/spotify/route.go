package spotify

import (
	"net/http"

	"github.com/ccrsxx/api/internal/features/auth"
)

type Config struct {
	Router         *http.ServeMux
	Service        *Service
	AuthMiddleware *auth.Middleware
}

func LoadRoutes(cfg Config) {
	mux := http.NewServeMux()

	ctrl := NewController(cfg.Service)

	mux.Handle("GET /currently-playing",
		cfg.AuthMiddleware.IsAuthorized(
			http.HandlerFunc(ctrl.getCurrentlyPlaying),
		),
	)

	cfg.Router.Handle("/spotify/", http.StripPrefix("/spotify", mux))
}
