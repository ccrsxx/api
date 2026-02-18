package logger

import (
	"log/slog"
	"os"

	"github.com/ccrsxx/api/internal/config"
)

func Init() {
	if config.Config().IsDevelopment {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		return
	}

	opts := slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				return slog.Attr{
					Key:   "severity",
					Value: a.Value,
				}
			}

			if a.Key == slog.MessageKey {
				return slog.Attr{
					Key:   "message",
					Value: a.Value,
				}
			}

			return a
		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))

	slog.SetDefault(logger)
}
