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

func LoadRoutes(cfg Config) {
	controller := NewController(mockIcon)

	cfg.Router.HandleFunc("GET /favicon.ico", controller.getFavicon)
}
