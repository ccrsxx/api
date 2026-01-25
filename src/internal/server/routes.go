package server

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/features/docs"
	"github.com/ccrsxx/api-go/src/internal/features/favicon"
	"github.com/ccrsxx/api-go/src/internal/features/home"
	"github.com/ccrsxx/api-go/src/internal/features/jellyfin"
	"github.com/ccrsxx/api-go/src/internal/features/spotify"
	"github.com/ccrsxx/api-go/src/internal/features/sse"
	"github.com/ccrsxx/api-go/src/internal/features/tools"
	"github.com/ccrsxx/api-go/src/internal/middleware"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	sse.LoadRoutes(router)
	home.LoadRoutes(router)
	docs.LoadRoutes(router)
	tools.LoadRoutes(router)
	favicon.LoadRoutes(router)
	spotify.LoadRoutes(router)
	jellyfin.LoadRoutes(router)

	routes := middleware.Recovery(
		middleware.Cors(
			middleware.Logging(
				middleware.RateLimit(100, 1*time.Minute)(
					router,
				),
			),
		),
	)

	return routes

}
