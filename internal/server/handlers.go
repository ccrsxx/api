package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/cache"
	cloudflareClient "github.com/ccrsxx/api/internal/clients/cloudflare"
	githubClient "github.com/ccrsxx/api/internal/clients/github"
	gmailClient "github.com/ccrsxx/api/internal/clients/gmail"
	ipInfoClient "github.com/ccrsxx/api/internal/clients/ipinfo"
	jellyfinClient "github.com/ccrsxx/api/internal/clients/jellyfin"
	navidromeClient "github.com/ccrsxx/api/internal/clients/navidrome"
	pixivClient "github.com/ccrsxx/api/internal/clients/pixiv"
	pushoverClient "github.com/ccrsxx/api/internal/clients/pushover"
	spotifyClient "github.com/ccrsxx/api/internal/clients/spotify"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/features/auth"
	"github.com/ccrsxx/api/internal/features/contacts"
	"github.com/ccrsxx/api/internal/features/contents"
	"github.com/ccrsxx/api/internal/features/docs"
	"github.com/ccrsxx/api/internal/features/favicon"
	"github.com/ccrsxx/api/internal/features/guestbook"
	"github.com/ccrsxx/api/internal/features/home"
	"github.com/ccrsxx/api/internal/features/jellyfin"
	"github.com/ccrsxx/api/internal/features/likes"
	"github.com/ccrsxx/api/internal/features/navidrome"
	"github.com/ccrsxx/api/internal/features/og"
	"github.com/ccrsxx/api/internal/features/pixiv"
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

	gmailClient := gmailClient.NewClient(gmailClient.Config{
		Username: cfg.EmailAddress,
		Password: cfg.EmailPassword,
	})

	cloudflareClient := cloudflareClient.NewClient(cloudflareClient.Config{
		SecretKey: cfg.CloudflareTurnstileSecretKey,
	})

	pixivClient := pixivClient.NewClient(pixivClient.Config{
		Token: cfg.PixivToken,
	})

	pushoverClient := pushoverClient.NewClient(pushoverClient.Config{
		UserKey:  cfg.PushoverUserKey,
		AppToken: cfg.PushoverAppToken,
	})

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

	navidromeClient := navidromeClient.NewClient(navidromeClient.Config{
		URL:      cfg.NavidromeURL,
		Username: cfg.NavidromeUsername,
		Password: cfg.NavidromePassword,
	})

	spotifyService := spotify.NewService(spotify.ServiceConfig{
		Client: spotifyClient,
	})

	jellyfinService := jellyfin.NewService(jellyfin.ServiceConfig{
		Client:           jellyfinClient,
		JellyfinUsername: cfg.JellyfinUsername,
		JellyfinImageURL: cfg.JellyfinImageURL,
	})

	navidromeService := navidrome.NewService(navidrome.ServiceConfig{
		Client:            navidromeClient,
		BackendPublicURL:  cfg.BackendPublicURL,
		NavidromeUsername: cfg.NavidromeUsername,
	})

	authService := auth.ServiceConfig{
		Pool:              pool,
		Database:          &auth.AuthDatabaseWrapper{Queries: db},
		SecretKey:         cfg.SecretKey,
		JwtSecret:         cfg.JWTSecret,
		GithubClient:      githubClient,
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
				IPInfoClient: ipInfoClient,
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
				AppContext:       ctx,
				SpotifyService:   spotifyService,
				JellyfinService:  jellyfinService,
				NavidromeService: navidromeService,
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

	pixiv.LoadRoutes(
		pixiv.Config{
			Router:         router,
			AuthMiddleware: privateAuthMiddleware,
			Service: pixiv.NewService(pixiv.ServiceConfig{
				Client:        pixivClient,
				PixivImageURL: cfg.PixivImageURL,
			}),
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

	navidrome.LoadRoutes(
		navidrome.Config{
			Router:         router,
			Service:        navidromeService,
			AuthMiddleware: publicAuthMiddleware,
		},
	)

	contacts.LoadRoutes(
		contacts.Config{
			Router:         router,
			AppContext:     ctx,
			AuthMiddleware: publicAuthMiddleware,
			Service: contacts.NewService(contacts.ServiceConfig{
				Database:         db,
				EmailClient:      gmailClient,
				EmailTarget:      cfg.EmailTarget,
				EmailAddress:     cfg.EmailAddress,
				PushoverClient:   pushoverClient,
				CloudflareClient: cloudflareClient,
			}),
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
				Database:     db,
				EmailClient:  gmailClient,
				EmailTarget:  cfg.EmailTarget,
				EmailAddress: cfg.EmailAddress,
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
