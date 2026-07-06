package contacts

import (
	"net/http"
)

type Config struct {
	Router  *http.ServeMux
	Service *Service
}

func LoadRoutes(cfg Config) {
	ctrl := NewController(cfg.Service)

	cfg.Router.HandleFunc("POST /contacts", ctrl.CreateContact)
}
