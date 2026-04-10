package config

import (
	"testing"
)

// Helper to quickly set all required env vars for a successful parse
func setValidEnv(t *testing.T, appEnv string) {
	t.Setenv("PORT", "8080")
	t.Setenv("OG_URL", "http://localhost")
	t.Setenv("APP_ENV", appEnv)
	t.Setenv("SECRET_KEY", "secret")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")
	t.Setenv("IPINFO_TOKEN", "token")
	t.Setenv("JELLYFIN_URL", "url")
	t.Setenv("JELLYFIN_API_KEY", "key")
	t.Setenv("JELLYFIN_USERNAME", "user")
	t.Setenv("JELLYFIN_IMAGE_URL", "img")
	t.Setenv("SPOTIFY_CLIENT_ID", "id")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "secret")
	t.Setenv("SPOTIFY_REFRESH_TOKEN", "refresh")
}

func TestLoad_Success(t *testing.T) {
	t.Run("Development Mode", func(t *testing.T) {
		setValidEnv(t, string(EnvironmentDevelopment))

		cfg := Init()

		if cfg.Port != 8080 {
			t.Errorf("got port %d, want 8080", cfg.Port)
		}

		if len(cfg.AllowedOrigins) != 2 {
			t.Errorf("got %d allowed origins, want 2", len(cfg.AllowedOrigins))
		}

		if !cfg.IsDevelopment {
			t.Error("want IsDevelopment to be true")
		}

		if cfg.IsProduction {
			t.Error("want IsProduction to be false")
		}
	})

	t.Run("Production Mode", func(t *testing.T) {
		setValidEnv(t, string(EnvironmentProduction))

		cfg := Init()

		if !cfg.IsProduction {
			t.Error("want IsProduction to be true")
		}

		if cfg.IsDevelopment {
			t.Error("want IsDevelopment to be false")
		}
	})
}

func TestLoad_PanicOnMissingEnv(t *testing.T) {
	// We only set PORT, causing env.Parse to fail on all the other required fields
	t.Setenv("PORT", "8080")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Load() did not panic on missing required env vars")
		}
	}()

	Init()
}

func TestEnvironmentApp_UnmarshalText(t *testing.T) {
	t.Run("Valid Environments", func(t *testing.T) {
		var env EnvironmentApp

		if err := env.UnmarshalText([]byte("development")); err != nil {
			t.Errorf("unwanted error: %v", err)
		}

		if env != EnvironmentDevelopment {
			t.Errorf("got %v, want development", env)
		}

		if err := env.UnmarshalText([]byte("production")); err != nil {
			t.Errorf("unwanted error: %v", err)
		}

		if env != EnvironmentProduction {
			t.Errorf("got %v, want production", env)
		}
	})

	t.Run("Invalid Environment", func(t *testing.T) {
		var env EnvironmentApp

		err := env.UnmarshalText([]byte("staging"))

		if err == nil {
			t.Error("want error for invalid environment")
		}

		if env != "" { // Ensure it wasn't set
			t.Errorf("got %v, want empty", env)
		}
	})
}

func TestLoad_PanicOnInvalidAppEnv(t *testing.T) {
	setValidEnv(t, "invalid_env_string")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Load() did not panic on invalid APP_ENV text")
		}
	}()

	Init()
}
