package views

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

	mux.HandleFunc("GET /{slug}", ctrl.GetViewCount)
	mux.HandleFunc("POST /{slug}", ctrl.IncrementView)

	cfg.Router.Handle("/views/", http.StripPrefix("/views", mux))
}
