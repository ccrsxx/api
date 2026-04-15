package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/cache"
	"github.com/ccrsxx/api/internal/clients/ipinfo"
	jellyfinClient "github.com/ccrsxx/api/internal/clients/jellyfin"
	spotifyClient "github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/docs"
	"github.com/ccrsxx/api/internal/features/favicon"
	"github.com/ccrsxx/api/internal/features/home"
	"github.com/ccrsxx/api/internal/features/jellyfin"
	"github.com/ccrsxx/api/internal/features/og"
	"github.com/ccrsxx/api/internal/features/spotify"
	"github.com/ccrsxx/api/internal/features/sse"
	"github.com/ccrsxx/api/internal/features/tools"
	"github.com/ccrsxx/api/internal/features/users"
	"github.com/ccrsxx/api/internal/middleware"
)

func LoadHandlers(ctx context.Context, cfg config.AppConfig, db *sqlc.Queries) http.Handler {
	router := http.NewServeMux()

	memoryCache := cache.NewMemoryCache(ctx, cache.DefaultCleanupInterval)

	ipInfoClient := ipinfo.NewClient(cfg.IPInfoToken)

	spotifyClient := spotifyClient.NewClient(spotifyClient.Config{
		ClientID:     cfg.SpotifyClientID,
		MemoryCache:  memoryCache,
		ClientSecret: cfg.SpotifyClientSecret,
		RefreshToken: cfg.SpotifyRefreshToken,
	})

	jellyfinClient := jellyfinClient.NewClient(jellyfinClient.Config{
		URL:      cfg.JellyfinURL,
		APIKey:   cfg.JellyfinAPIKey,
		ImageURL: cfg.JellyfinImageURL,
		Username: cfg.JellyfinUsername,
	})

	usersServices := users.NewService(users.ServiceConfig{
		Database: db,
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
		JellyfinImageURL: cfg.JellyfinImageURL,
	})

	toolsController := tools.NewController(
		tools.NewService(
			tools.ServiceConfig{
				Fetcher: ipInfoClient.GetIPInfo,
			},
		),
	)

	// Shared rate-limited handler for GetIpInfo. Limits to 10 requests per 10 seconds.
	sharedGetIPInfoController := middleware.RateLimit(ctx, 10, 10*time.Second)(
		http.HandlerFunc(toolsController.GetIPInfo),
	)

	og.LoadRoutes(og.Config{
		Router: router,
		Service: og.NewService(og.ServiceConfig{
			OgURL: cfg.OgURL,
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
			AppContext:     ctx,
			AuthMiddleware: authMiddleware,
		},
	)

	home.LoadRoutes(
		home.Config{
			Router:                    router,
			ToolsController:           toolsController,
			SharedGetIPInfoController: sharedGetIPInfoController,
		},
	)

	docs.LoadRoutes(
		docs.Config{
			Router: router,
		},
	)

	users.LoadRoutes(
		users.Config{
			Router:  router,
			Service: usersServices,
		},
	)

	tools.LoadRoutes(
		tools.Config{
			Router:                    router,
			ToolsController:           toolsController,
			SharedGetIPInfoController: sharedGetIPInfoController,
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
				middleware.RateLimit(ctx, 100, 1*time.Minute)(
					router,
				),
			),
		),
	)

	return handlers
}
