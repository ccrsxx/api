package server

import (
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/cache"
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

func LoadHandlers(cfg config.AppConfig) http.Handler {
	router := http.NewServeMux()

	memoryCache := cache.NewMemoryCache(cache.DefaultCleanupInterval)

	ipInfoClient := ipinfo.NewClient(cfg.IpInfoToken)

	spotifyClient := spotifyClient.NewClient(spotifyClient.Config{
		ClientID:     cfg.SpotifyClientID,
		MemoryCache:  memoryCache,
		ClientSecret: cfg.SpotifyClientSecret,
		RefreshToken: cfg.SpotifyRefreshToken,
	})

	jellyfinClient := jellyfinClient.NewClient(jellyfinClient.Config{
		URL:      cfg.JellyfinUrl,
		ApiKey:   cfg.JellyfinApiKey,
		ImageURL: cfg.JellyfinImageUrl,
		Username: cfg.JellyfinUsername,
	})

	authMiddleware := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{
		SecretKey: cfg.SecretKey,
	}))

	spotifyService := spotify.NewService(spotify.ServiceConfig{
		Fetcher: spotifyClient.GetCurrentlyPlaying,
	})

	jellyfinService := jellyfin.NewService(jellyfin.ServiceConfig{
		Fetcher:          jellyfinClient.GetSessions,
		JellyfinUsername: cfg.JellyfinUsername,
		JellyfinImageUrl: cfg.JellyfinImageUrl,
	})

	toolsController := tools.NewController(
		tools.NewService(
			tools.ServiceConfig{
				Fetcher: ipInfoClient.GetIPInfo,
			},
		),
	)

	// Shared rate-limited handler for GetIpInfo. Limits to 10 requests per 10 seconds.
	sharedGetIpInfoController := middleware.RateLimit(10, 10*time.Second)(
		http.HandlerFunc(toolsController.GetIpInfo),
	)

	og.LoadRoutes(og.Config{
		Router: router,
		Service: og.NewService(og.ServiceConfig{
			OgUrl: cfg.OgUrl,
		}),
		ControllerConfig: og.ControllerConfig{
			IsProduction: cfg.IsProduction,
		},
	})

	sse.LoadRoutes(
		sse.Config{
			Router: router,
			Service: sse.NewService(sse.ServiceConfig{
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

	handlers := middleware.Recovery(
		middleware.Cors(cfg.AllowedOrigins)(
			middleware.Logging(
				middleware.RateLimit(100, 1*time.Minute)(
					router,
				),
			),
		),
	)

	return handlers
}
