package og

import (
	"net/http"
)

type Config struct {
	ControllerConfig ControllerConfig
}

func LoadRoutes(router *http.ServeMux, service *Service, config Config) {
	controller := NewController(service, config)

	router.HandleFunc("/og", controller.getOg)
}
