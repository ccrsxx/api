package metrics

import (
	"net/http"
)

type Config struct {
	Router *http.ServeMux
}

func LoadRoutes(cfg Config) {
	ctrl := NewController()

	cfg.Router.HandleFunc("GET /metrics", ctrl.GetMetrics)
}
