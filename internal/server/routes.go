package server

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/features/docs"
	"github.com/ccrsxx/api/internal/features/favicon"
	"github.com/ccrsxx/api/internal/features/home"
	"github.com/ccrsxx/api/internal/features/jellyfin"
	"github.com/ccrsxx/api/internal/features/og"
	"github.com/ccrsxx/api/internal/features/spotify"
	"github.com/ccrsxx/api/internal/features/sse"
	"github.com/ccrsxx/api/internal/features/tools"
	"github.com/ccrsxx/api/internal/middleware"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	svcSse := sse.NewService(
		sse.Config{
			PollInterval:    1 * time.Second,
			SpotifyFetcher:  spotify.Service.GetCurrentlyPlaying,
			JellyfinFetcher: jellyfin.Service.GetCurrentlyPlaying,
		},
	)

	sse.LoadRoutes(router, svcSse)

	og.LoadRoutes(router)
	home.LoadRoutes(router)
	docs.LoadRoutes(router)
	tools.LoadRoutes(router)
	favicon.LoadRoutes(router)
	spotify.LoadRoutes(router)
	jellyfin.LoadRoutes(router)

	routes := middleware.Recovery(
		middleware.Cors(config.Env().AllowedOrigins)(
			middleware.Logging(
				middleware.RateLimit(100, 1*time.Minute)(
					router,
				),
			),
		),
	)

	return routes
}
