package docs

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.json
var openapiSpec []byte

type Config struct {
	Router *http.ServeMux
}

func LoadRoutes(config Config) {
	controller := NewController(openapiSpec)

	config.Router.HandleFunc("GET /docs", controller.getDocs)
}
