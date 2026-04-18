package likes

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

	mux.HandleFunc("GET /{slug}", ctrl.GetLikeStatus)

	mux.Handle("POST /{slug}",
		cfg.AuthMiddleware.IsAuthorized(
			http.HandlerFunc(ctrl.IncrementLike),
		),
	)

	cfg.Router.Handle("/likes/", http.StripPrefix("/likes", mux))
}
