package guestbook

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

	mux.HandleFunc("GET /", ctrl.GetGuestbook)

	mux.Handle("POST /",
		cfg.AuthMiddleware.IsAuthorizedFromOauth(
			http.HandlerFunc(ctrl.CreateGuestbook),
		),
	)

	mux.Handle("DELETE /{id}",
		cfg.AuthMiddleware.IsAuthorizedFromOauth(
			cfg.AuthMiddleware.IsAdminFromOauth(
				http.HandlerFunc(ctrl.DeleteGuestbook),
			),
		),
	)

	cfg.Router.Handle("/guestbook/", http.StripPrefix("/guestbook", mux))
}
