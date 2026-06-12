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

func LoadRoutes(cfg Config) {
	ctrl := NewController(openapiSpec)

	cfg.Router.HandleFunc("GET /docs", ctrl.GetDocs)
}
