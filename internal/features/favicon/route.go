package favicon

import (
	_ "embed"
	"net/http"
)

//go:embed favicon.ico
var mockIcon []byte

func LoadRoutes(router *http.ServeMux) {
	controller := NewController(mockIcon)

	router.HandleFunc("GET /favicon.ico", controller.getFavicon)
}
