package og

import (
	"net/http"
)

type Config struct {
	Router           *http.ServeMux
	Service          *Service
	ControllerConfig ControllerConfig
}

func LoadRoutes(cfg Config) {
	ctrl := NewController(cfg.Service, cfg.ControllerConfig)

	cfg.Router.HandleFunc("/og", ctrl.getOg)
}
