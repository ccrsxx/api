package sse

import (
	"net/http"

	"github.com/ccrsxx/api-go/src/internal/features/auth"
)

func LoadRoutes(router *http.ServeMux) {
	router.Handle("GET /sse",
		auth.Middleware.IsAuthorizedFromQuery(
			Middleware.IsConnectionAllowed(
				http.HandlerFunc(Controller.getCurrentPlayingSSE),
			),
		),
	)

}
