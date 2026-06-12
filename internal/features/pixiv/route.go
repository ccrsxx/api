package pixiv

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

	mux.Handle("GET /bookmarks",
		cfg.AuthMiddleware.IsAuthorizedFromBearer(
			http.HandlerFunc(ctrl.GetBookmarks),
		),
	)

	mux.Handle("GET /bookmarks/all",
		cfg.AuthMiddleware.IsAuthorizedFromBearer(
			http.HandlerFunc(ctrl.GetAllBookmarks),
		),
	)

	cfg.Router.Handle("/pixiv/", http.StripPrefix("/pixiv", mux))
}
