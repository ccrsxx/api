package auth

import (
	"net/http"
)

type Config struct {
	Router         *http.ServeMux
	Service        *Service
	AuthMiddleware *Middleware
}

func LoadRoutes(cfg Config) {
	mux := http.NewServeMux()

	ctrl := NewController(cfg.Service)

	mux.Handle("GET /me",
		cfg.AuthMiddleware.IsAuthorizedFromOauth(
			http.HandlerFunc(ctrl.GetCurrentUser),
		),
	)

	mux.HandleFunc("GET /github/login", ctrl.LoginGithub)

	mux.HandleFunc("POST /github/logout", ctrl.LogoutGithub)

	mux.HandleFunc("GET /github/callback", ctrl.LoginGithubCallback)

	cfg.Router.Handle("/auth/", http.StripPrefix("/auth", mux))
}
