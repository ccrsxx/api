package og

import (
	"net/http"
)

type Config struct {
	Router           *http.ServeMux
	Service          *Service
	ControllerConfig ControllerConfig
}

func LoadRoutes(config Config) {
	controller := NewController(config.Service, config.ControllerConfig)

	config.Router.HandleFunc("/og", controller.getOg)
}
