package config

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvProduction  Environment = "production"
	EnvDevelopment Environment = "development"
)

func (e *Environment) UnmarshalText(text []byte) error {
	val := Environment(text)

	switch val {
	case EnvDevelopment, EnvProduction:
		*e = val
		return nil
	default:
		return fmt.Errorf("invalid app env: %s", val)
	}
}

type appEnv struct {
	Port           string      `env:"PORT,required"`
	AppEnv         Environment `env:"APP_ENV,required"`
	SecretKey      string      `env:"SECRET_KEY,required"`
	IpInfoToken    string      `env:"IPINFO_TOKEN,required"`
	AllowedOrigins []string    `env:"ALLOWED_ORIGINS,required"`
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
		// 1. Try to load .env.local first (Dev/Overrides)
		// We ignore errors because in Production, this file won't exist.
		_ = godotenv.Load(".env.local")

		// 2. Try to load .env (Defaults)
		// If .env.local already set a var, this line WON'T overwrite it.
		// If on Docker (Production), these might fail but System Envs will take over.
		_ = godotenv.Load(".env")

		if err := env.Parse(&envInstance); err != nil {
			slog.Error("env parse error", "error", err)
			os.Exit(1)
		}
	})
}
