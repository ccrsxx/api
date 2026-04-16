package contents

import (
	"net/http"
)

type Config struct {
	Router  *http.ServeMux
	Service *Service
}

func LoadRoutes(cfg Config) {
	mux := http.NewServeMux()

	ctrl := NewController(cfg.Service)

	mux.HandleFunc("GET /{type}", ctrl.GetContentData)

	mux.HandleFunc("POST /{type}", ctrl.UpsertContent)

	cfg.Router.Handle("/content/", http.StripPrefix("/content", mux))
}
