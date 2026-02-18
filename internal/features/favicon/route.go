package favicon

import (
	"net/http"
)

func LoadRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /favicon.ico", Controller.getFavicon)
}
