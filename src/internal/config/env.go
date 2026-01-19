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

		slog.Info("env loading", "file", envFile)

		if err := godotenv.Load(envFile); err != nil {
			if Config().IsDevelopment {
				slog.Error("env load error", "file", envFile, "error", err)
				os.Exit(1)
			}

			slog.Info("env file missing", "action", "using system env")
		}

		if err := env.Parse(&envInstance); err != nil {
			slog.Error("env parse error", "error", err)
			os.Exit(1)
		}

		slog.Info("env loaded")
	})
}
