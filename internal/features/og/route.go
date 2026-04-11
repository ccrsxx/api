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
	controller := NewController(cfg.Service, cfg.ControllerConfig)

	cfg.Router.HandleFunc("/og", controller.getOg)
}
