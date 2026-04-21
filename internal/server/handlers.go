package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/cache"
	githubClient "github.com/ccrsxx/api/internal/clients/github"
	ipInfoClient "github.com/ccrsxx/api/internal/clients/ipinfo"
	jellyfinClient "github.com/ccrsxx/api/internal/clients/jellyfin"
	spotifyClient "github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/contents"
	"github.com/ccrsxx/api/internal/features/docs"
	"github.com/ccrsxx/api/internal/features/favicon"
	"github.com/ccrsxx/api/internal/features/guestbook"
	"github.com/ccrsxx/api/internal/features/home"
	"github.com/ccrsxx/api/internal/features/jellyfin"
	"github.com/ccrsxx/api/internal/features/likes"
	"github.com/ccrsxx/api/internal/features/og"
	"github.com/ccrsxx/api/internal/features/spotify"
	"github.com/ccrsxx/api/internal/features/sse"
	"github.com/ccrsxx/api/internal/features/statistics"
	"github.com/ccrsxx/api/internal/features/tools"
	"github.com/ccrsxx/api/internal/features/views"
	"github.com/ccrsxx/api/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func LoadHandlers(ctx context.Context, cfg config.AppConfig, pool *pgxpool.Pool, db *sqlc.Queries) http.Handler {
	router := http.NewServeMux()

	memoryCache := cache.NewMemoryCache(ctx, cache.DefaultCleanupInterval)

	ipInfoClient := ipInfoClient.NewClient(cfg.IPInfoToken)

	githubClient := githubClient.NewClient(githubClient.Config{})

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

	spotifyService := spotify.NewService(spotify.ServiceConfig{
		Fetcher: spotifyClient.GetCurrentlyPlaying,
	})

	jellyfinService := jellyfin.NewService(jellyfin.ServiceConfig{
		Fetcher:          jellyfinClient.GetSessions,
		JellyfinUsername: cfg.JellyfinUsername,
		JellyfinImageURL: cfg.JellyfinImageURL,
	})

	authService := auth.ServiceConfig{
		Pool:              pool,
		Database:          &auth.AuthDatabaseWrapper{Queries: db},
		SecretKey:         cfg.SecretKey,
		JwtSecret:         cfg.JWTSecret,
		GetGithubUser:     githubClient.GetCurrentUser,
		FrontendBaseURL:   cfg.FrontendBaseURL,
		FrontendPublicURL: cfg.FrontendPublicURL,
		GithubOauthConfig: &oauth2.Config{
			Endpoint:     github.Endpoint,
			ClientID:     cfg.OauthGithubClientID,
			ClientSecret: cfg.OauthGithubClientSecret,
			Scopes:       []string{"read:user"},
		},
	}

	publicAuthMiddleware := auth.NewMiddleware(auth.NewService(authService))

	privateAuthMiddleware := auth.NewMiddleware(auth.NewService(auth.ServiceConfig{
		SecretKey: cfg.PrivateSecretKey,
	}))

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

	auth.LoadRoutes(
		auth.Config{
			Router:         router,
			Service:        auth.NewService(authService),
			AuthMiddleware: publicAuthMiddleware,
		},
	)

	og.LoadRoutes(
		og.Config{
			Router: router,
			Service: og.NewService(og.ServiceConfig{
				OgURL: cfg.OgURL,
			}),
			ControllerConfig: og.ControllerConfig{
				IsProduction: cfg.IsProduction,
			},
		},
	)

	sse.LoadRoutes(
		sse.Config{
			Router: router,
			Service: sse.NewService(sse.ServiceConfig{
				SpotifyFetcher:  spotifyService.GetCurrentlyPlaying,
				JellyfinFetcher: jellyfinService.GetCurrentlyPlaying,
			}),
			AppContext:     ctx,
			AuthMiddleware: publicAuthMiddleware,
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

	views.LoadRoutes(
		views.Config{
			Router: router,
			Service: views.NewService(views.ServiceConfig{
				Database: db,
			}),
			AppContext:     ctx,
			AuthMiddleware: publicAuthMiddleware,
		},
	)

	likes.LoadRoutes(
		likes.Config{
			Router: router,
			Service: likes.NewService(likes.ServiceConfig{
				Database: db,
			}),
			AuthMiddleware: publicAuthMiddleware,
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
			AuthMiddleware: publicAuthMiddleware,
		},
	)

	jellyfin.LoadRoutes(
		jellyfin.Config{
			Router:         router,
			Service:        jellyfinService,
			AuthMiddleware: publicAuthMiddleware,
		},
	)

	contents.LoadRoutes(
		contents.Config{
			Router: router,
			Service: contents.NewService(contents.ServiceConfig{
				Database: db,
			}),
			AuthMiddleware: privateAuthMiddleware,
		},
	)

	statistics.LoadRoutes(
		statistics.Config{
			Router: router,
			Service: statistics.NewService(statistics.ServiceConfig{
				Database: db,
			}),
		},
	)

	guestbook.LoadRoutes(
		guestbook.Config{
			Router:         router,
			AuthMiddleware: publicAuthMiddleware,
			Service: guestbook.NewService(guestbook.ServiceConfig{
				Database: db,
			}),
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
