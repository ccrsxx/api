package docs

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.json
var openapiSpec []byte

func LoadRoutes(router *http.ServeMux) {
	controller := NewController(openapiSpec)

	router.HandleFunc("GET /docs", controller.getDocs)
}
