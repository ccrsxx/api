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

	controller := NewController(cfg.Service)

	mux.Handle("GET /currently-playing",
		cfg.AuthMiddleware.IsAuthorized(
			http.HandlerFunc(controller.getCurrentlyPlaying),
		),
	)

	cfg.Router.Handle("/spotify/", http.StripPrefix("/spotify", mux))
}
