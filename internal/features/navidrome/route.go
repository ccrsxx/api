package navidrome

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
		cfg.AuthMiddleware.IsAuthorizedFromBearer(
			http.HandlerFunc(ctrl.GetCurrentlyPlaying),
		),
	)

	mux.HandleFunc("GET /cover-art/{coverArtID}", ctrl.GetCoverArt)

	cfg.Router.Handle("/navidrome/", http.StripPrefix("/navidrome", mux))
}
