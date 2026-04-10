package spotify

import (
	"net/http"

	"github.com/ccrsxx/api/internal/features/auth"
)

func LoadRoutes(router *http.ServeMux, service *service) {
	mux := http.NewServeMux()

	controller := NewController(service)

	mux.Handle("GET /currently-playing",
		auth.Middleware.IsAuthorized(
			http.HandlerFunc(controller.getCurrentlyPlaying),
		),
	)

	router.Handle("/spotify/", http.StripPrefix("/spotify", mux))
}
