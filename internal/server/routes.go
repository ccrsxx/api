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

	authService := auth.NewService(auth.ServiceConfig{
		SecretKey: config.Env().SecretKey,
	})

	authMiddleware := auth.NewMiddleware(authService)

	configOg := og.Config{ControllerConfig: og.ControllerConfig{
		IsProduction: config.Config().IsProduction,
	}}

	serviceOg := og.NewService(og.ServiceConfig{
		OgUrl:      config.Env().OgUrl,
		HttpClient: &http.Client{Timeout: 8 * time.Second},
	})

	serviceSpotify := spotify.NewService(spotify.ServiceConfig{
		Fetcher: spotifyClient.DefaultClient().GetCurrentlyPlaying,
	})

	serviceJellyfin := jellyfin.NewService(jellyfin.ServiceConfig{
		Fetcher:          jellyfinClient.DefaultClient().GetSessions,
		JellyfinUsername: config.Env().JellyfinUsername,
	})

	serviceSse := sse.NewService(
		sse.ServiceConfig{
			PollInterval:    1 * time.Second,
			SpotifyFetcher:  serviceSpotify.GetCurrentlyPlaying,
			JellyfinFetcher: serviceJellyfin.GetCurrentlyPlaying,
		},
	)

	serviceTools := tools.NewService(
		tools.ServiceConfig{
			Fetcher: ipinfo.DefaultClient().GetIPInfo,
		},
	)

	controllerTools := tools.NewController(serviceTools)

	configTools := tools.Config{
		ToolsController: controllerTools,
		SharedGetIpInfo: http.HandlerFunc(controllerTools.GetIpInfo),
	}

	og.LoadRoutes(router, serviceOg, configOg)
	sse.LoadRoutes(router, serviceSse, authMiddleware)
	home.LoadRoutes(router, configTools)
	docs.LoadRoutes(router)
	tools.LoadRoutes(router, configTools)
	favicon.LoadRoutes(router)
	spotify.LoadRoutes(router, serviceSpotify, authMiddleware)
	jellyfin.LoadRoutes(router, serviceJellyfin, authMiddleware)

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
