package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appEnv struct {
	Port      string `env:"PORT,required"`
	SecretKey string `env:"SECRET_KEY,required"`
}

var envInstance appEnv

var once sync.Once

func Env() *appEnv {
	return &envInstance
}

func LoadEnv() {
	once.Do(func() {
		envFile := ".env"

		if Config().IsDevelopment {
			envFile = ".env.local"
		}

		log.Println("Loading environment variables from", envFile)

		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("Failed to load %s file: %v", envFile, err)
		}

		if err := env.Parse(&envInstance); err != nil {
			log.Fatalf("Failed to parse env vars: %v", err)
		}

		log.Printf("Environment variables loaded successfully")
	})
}
