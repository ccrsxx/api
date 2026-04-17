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

	mux.HandleFunc("GET /", ctrl.GetContentData)

	mux.HandleFunc("POST /", ctrl.UpsertContent)

	cfg.Router.Handle("/contents/", http.StripPrefix("/contents", mux))
}
