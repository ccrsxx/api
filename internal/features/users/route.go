package users

import (
	"net/http"
)

type Config struct {
	Router  *http.ServeMux
	Service *Service
}

func LoadRoutes(cfg Config) {
	mux := http.NewServeMux()

	ctrl := NewController(cfg.Service)

	mux.HandleFunc("GET /", ctrl.GetListUsers)

	cfg.Router.Handle("/users/", http.StripPrefix("/users", mux))
}
