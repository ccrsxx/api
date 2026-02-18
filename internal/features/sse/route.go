package sse

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/features/auth"
	m "github.com/ccrsxx/api/internal/middleware"
)

func LoadRoutes(router *http.ServeMux) {
	router.Handle("GET /sse",
		auth.Middleware.IsAuthorizedFromQuery(
			m.RateLimit(10, 10*time.Second)(
				Middleware.IsConnectionAllowed(
					http.HandlerFunc(Controller.getCurrentPlayingSSE),
				),
			),
		),
	)
}
