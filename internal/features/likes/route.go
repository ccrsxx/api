package likes

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

	mux.HandleFunc("GET /{slug}", ctrl.GetLikeStatus)

	mux.HandleFunc("POST /{slug}", ctrl.IncrementLike)

	cfg.Router.Handle("/likes/", http.StripPrefix("/likes", mux))
}
