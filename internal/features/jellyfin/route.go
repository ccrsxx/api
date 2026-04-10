package jellyfin

import (
	"net/http"

	"github.com/ccrsxx/api/internal/features/auth"
)

func LoadRoutes(router *http.ServeMux, service *Service, authMiddleware *auth.Middleware) {
	mux := http.NewServeMux()

	controller := NewController(service)

	mux.Handle("GET /currently-playing",
		authMiddleware.IsAuthorized(
			http.HandlerFunc(controller.getCurrentlyPlaying),
		),
	)

	router.Handle("/jellyfin/", http.StripPrefix("/jellyfin", mux))
}
