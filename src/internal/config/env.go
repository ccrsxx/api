package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appEnv struct {
	Port           string   `env:"PORT,required"`
	SecretKey      string   `env:"SECRET_KEY,required"`
	IpInfoToken    string   `env:"IPINFO_TOKEN,required"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS,required"`
}

var (
	once        sync.Once
	envInstance appEnv
)

func Env() *appEnv {
	return &envInstance
}

func LoadEnv() {
	once.Do(func() {
		envFile := ".env"

		if Config().IsDevelopment {
			envFile = ".env.local"
		}

		slog.Info("Loading environment variables", "file", envFile)

		if err := godotenv.Load(envFile); err != nil {
			if Config().IsDevelopment {
				slog.Error("Failed to load env file", "file", envFile, "error", err)
				os.Exit(1)
			}

			slog.Info("No env file found, proceeding with system environment variables", "file", envFile)
		}

		if err := env.Parse(&envInstance); err != nil {
			slog.Error("Failed to parse env vars", "error", err)
		}

		slog.Info("Environment variables loaded successfully")
	})
}
