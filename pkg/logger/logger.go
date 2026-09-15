package logger

import (
	"log/slog"
	"os"
)

// InitLogger initializes and sets the default structured JSON logger.
func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
			}
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}
