package server

import (
	"net/http"
	"time"

	jellyfinClient "github.com/ccrsxx/api/internal/clients/jellyfin"
	spotifyClient "github.com/ccrsxx/api/internal/clients/spotify"
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

	serviceSpotify := spotify.NewService(spotify.Config{
		Fetcher: spotifyClient.DefaultClient().GetCurrentlyPlaying,
	})

	serviceJellyfin := jellyfin.NewService(jellyfin.Config{
		Fetcher:          jellyfinClient.DefaultClient().GetSessions,
		JellyfinUsername: config.Env().JellyfinUsername,
	})

	serviceSse := sse.NewService(
		sse.Config{
			PollInterval:    1 * time.Second,
			SpotifyFetcher:  serviceSpotify.GetCurrentlyPlaying,
			JellyfinFetcher: serviceJellyfin.GetCurrentlyPlaying,
		},
	)

	og.LoadRoutes(router)
	sse.LoadRoutes(router, serviceSse)
	home.LoadRoutes(router)
	docs.LoadRoutes(router)
	tools.LoadRoutes(router)
	favicon.LoadRoutes(router)
	spotify.LoadRoutes(router, serviceSpotify)
	jellyfin.LoadRoutes(router, serviceJellyfin)

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
