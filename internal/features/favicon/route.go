package favicon

import (
	_ "embed"
	"net/http"
)

//go:embed favicon.ico
var mockIcon []byte

type Config struct {
	Router *http.ServeMux
}

func LoadRoutes(config Config) {
	controller := NewController(mockIcon)

	config.Router.HandleFunc("GET /favicon.ico", controller.getFavicon)
}
