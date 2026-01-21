package spotify

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/features/auth"
)

func LoadRoutes(router *http.ServeMux) {
	mux := http.NewServeMux()

	mux.Handle("GET /currently-playing",
		auth.Middleware.IsAuthorized(
			http.HandlerFunc(Controller.getCurrentlyPlaying),
		),
	)

	router.Handle("/spotify/", http.StripPrefix("/spotify", mux))
}
