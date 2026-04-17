package contents

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

	mux.HandleFunc("GET /", ctrl.GetContentData)

	mux.Handle("POST /",
		cfg.AuthMiddleware.IsAuthorized(
			http.HandlerFunc(ctrl.UpsertContent),
		),
	)

	cfg.Router.Handle("/contents/", http.StripPrefix("/contents", mux))
}
