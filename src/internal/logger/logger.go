package logger

import (
	"log/slog"
	"os"

	"github.com/ccrsxx/api-go/src/internal/config"
)

func Init() {
	if config.Config().IsDevelopment {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		return
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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
	})

	slog.SetDefault(slog.New(handler))

}
