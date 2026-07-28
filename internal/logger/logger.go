package logger

import (
	"log/slog"
	"os"

	"github.com/ccrsxx/api/internal/config"
)

// Load configures the global slog logger.
//
// Format (json|text) is decoupled from the environment so the JSON pipeline
// can be exercised locally. Level still follows the environment: debug in
// development, info in production.
//
// Every log line carries service_name and deployment_environment as base
// attributes so downstream collectors (Loki/Alloy) can index on them without
// the application having to repeat them at every call site.
//
// Note: we deliberately keep slog's default keys (level, msg, time). The
// previous GCP Cloud Logging renaming (severity/message) is gone because
// Grafana/Loki level detection and Logs Drilldown expect the standard keys.
func Load(cfg config.AppConfig) {
	level := slog.LevelInfo

	if cfg.IsDevelopment {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	if cfg.LogFormat == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	// Attach the context-aware handler so trace_id/request_id from the request
	// context are appended to every log line automatically.
	handler = NewContextHandler(handler)

	logger := slog.New(handler).With(
		slog.String("service_name", cfg.ServiceName),
		slog.String("deployment_environment", string(cfg.AppEnv)),
	)

	slog.SetDefault(logger)
}
