package sse

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/features/auth"
	middlewarePackage "github.com/ccrsxx/api-go/src/internal/middleware"
)

func LoadRoutes(router *http.ServeMux) {
	router.Handle("GET /sse",
		auth.Middleware.IsAuthorizedFromQuery(
			middlewarePackage.RateLimit(10, 10*time.Second)(
				Middleware.IsConnectionAllowed(
					http.HandlerFunc(Controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
