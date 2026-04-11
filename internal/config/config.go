package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type EnvironmentApp string

const (
	EnvironmentProduction  EnvironmentApp = "production"
	EnvironmentDevelopment EnvironmentApp = "development"
)

func (e *EnvironmentApp) UnmarshalText(text []byte) error {
	val := EnvironmentApp(text)

	switch val {
	case EnvironmentDevelopment, EnvironmentProduction:
		*e = val
		return nil
	default:
		return fmt.Errorf("invalid app env: %s", val)
	}
}

type AppConfig struct {
	Port           int            `env:"PORT,required"`
	OgUrl          string         `env:"OG_URL,required"`
	AppEnv         EnvironmentApp `env:"APP_ENV,required"`
	SecretKey      string         `env:"SECRET_KEY,required"`
	AllowedOrigins []string       `env:"ALLOWED_ORIGINS,required"`

	IpInfoToken string `env:"IPINFO_TOKEN,required"`

	JellyfinUrl      string `env:"JELLYFIN_URL,required"`
	JellyfinApiKey   string `env:"JELLYFIN_API_KEY,required"`
	JellyfinUsername string `env:"JELLYFIN_USERNAME,required"`
	JellyfinImageUrl string `env:"JELLYFIN_IMAGE_URL,required"`

	SpotifyClientID     string `env:"SPOTIFY_CLIENT_ID,required"`
	SpotifyClientSecret string `env:"SPOTIFY_CLIENT_SECRET,required"`
	SpotifyRefreshToken string `env:"SPOTIFY_REFRESH_TOKEN,required"`

	// Computed fields for convenience
	IsProduction  bool
	IsDevelopment bool
}

func Load() AppConfig {
	// System Environment Variables have the highest priority.
	// They override any loaded .env files.

	// Must load each .env files separately.
	// If we use godotenv.Load(".env.local", ".env"), it won't load .env if .env.local is missing.

	// 1. Try to load .env.local first (Dev/Overrides)
	// Use case: running development locally.
	// We ignore errors because in Production, this file won't exist.
	_ = godotenv.Load(".env.local")

	// 2. Try to load .env (Defaults)
	// Use case: running production locally without Docker.
	// If on Docker (Production), these might fail but System Envs will take over.
	_ = godotenv.Load(".env")

	var appConfig AppConfig

	// 3. Parse & Validate (The final check)
	// This reads from the actual environment (System + Loaded Files).
	// If "APP_ENV" is invalid or "PORT" is missing, this crashes the app HERE.

	if err := env.Parse(&appConfig); err != nil {
		slog.Error("env parse error", "error", err)

		// Panic instead of os.Exit(1) so we can test the failure state!
		panic(err)
	}

	appConfig.IsProduction = appConfig.AppEnv == EnvironmentProduction
	appConfig.IsDevelopment = appConfig.AppEnv == EnvironmentDevelopment

	return appConfig
}
