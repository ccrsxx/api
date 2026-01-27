package sse

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/src/internal/features/auth"
	middlewarepackage "github.com/ccrsxx/api/src/internal/middleware"
)

func LoadRoutes(router *http.ServeMux) {
	router.Handle("GET /sse",
		auth.Middleware.IsAuthorizedFromQuery(
			middlewarepackage.RateLimit(10, 10*time.Second)(
				Middleware.IsConnectionAllowed(
					http.HandlerFunc(Controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
