package jellyfin

import (
	"net/http"

	"github.com/ccrsxx/api/internal/features/auth"
)

type Config struct {
	Router         *http.ServeMux
	Service        *Service
	AuthMiddleware *auth.Middleware
}

func LoadRoutes(config Config) {
	mux := http.NewServeMux()

	controller := NewController(config.Service)

	mux.Handle("GET /currently-playing",
		config.AuthMiddleware.IsAuthorized(
			http.HandlerFunc(controller.getCurrentlyPlaying),
		),
	)

	config.Router.Handle("/jellyfin/", http.StripPrefix("/jellyfin", mux))
}
