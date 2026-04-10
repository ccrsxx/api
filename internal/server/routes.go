package server

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/clients/ipinfo"
	jellyfinClient "github.com/ccrsxx/api/internal/clients/jellyfin"
	spotifyClient "github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/features/auth"
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

	authMiddleware := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{
		SecretKey: config.Env().SecretKey,
	}))

	spotifyService := spotify.NewService(spotify.ServiceConfig{
		Fetcher: spotifyClient.DefaultClient().GetCurrentlyPlaying,
	})

	jellyfinService := jellyfin.NewService(jellyfin.ServiceConfig{
		Fetcher:          jellyfinClient.DefaultClient().GetSessions,
		JellyfinUsername: config.Env().JellyfinUsername,
	})

	toolsController := tools.NewController(
		tools.NewService(
			tools.ServiceConfig{
				Fetcher: ipinfo.DefaultClient().GetIPInfo,
			},
		),
	)

	sharedGetIpInfoController := http.HandlerFunc(toolsController.GetIpInfo)

	og.LoadRoutes(og.Config{
		Router: router,
		Service: og.NewService(og.ServiceConfig{
			OgUrl:      config.Env().OgUrl,
			HttpClient: &http.Client{Timeout: 8 * time.Second},
		}),
		ControllerConfig: og.ControllerConfig{IsProduction: config.Config().IsProduction},
	})

	sse.LoadRoutes(
		sse.Config{
			Router: router,
			Service: sse.NewService(sse.ServiceConfig{
				PollInterval:    1 * time.Second,
				SpotifyFetcher:  spotifyService.GetCurrentlyPlaying,
				JellyfinFetcher: jellyfinService.GetCurrentlyPlaying,
			}),
			AuthMiddleware: authMiddleware,
		},
	)

	home.LoadRoutes(
		home.Config{
			Router:                    router,
			ToolsController:           toolsController,
			SharedGetIpInfoController: sharedGetIpInfoController,
		},
	)

	docs.LoadRoutes(
		docs.Config{
			Router: router,
		},
	)

	tools.LoadRoutes(
		tools.Config{
			Router:                    router,
			ToolsController:           toolsController,
			SharedGetIpInfoController: sharedGetIpInfoController,
		},
	)

	favicon.LoadRoutes(
		favicon.Config{
			Router: router,
		},
	)

	spotify.LoadRoutes(
		spotify.Config{
			Router:         router,
			Service:        spotifyService,
			AuthMiddleware: authMiddleware,
		},
	)

	jellyfin.LoadRoutes(
		jellyfin.Config{
			Router:         router,
			Service:        jellyfinService,
			AuthMiddleware: authMiddleware,
		},
	)

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
