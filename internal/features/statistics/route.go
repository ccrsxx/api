package statistics

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

	mux.HandleFunc("GET /", ctrl.GetContentStatistics)

	cfg.Router.Handle("/statistics/", http.StripPrefix("/statistics", mux))
}
