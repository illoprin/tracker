package logger

import (
	"log/slog"
	"os"
	"tracker-backend/internal/config"
	"tracker-backend/internal/pkg/logger"
)

func InitLogger() (error, *slog.Logger) {
	var slogInstance *slog.Logger
	switch os.Getenv(config.EnvironmentEnvName) {
	case "local":
		opts := logger.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := logger.NewPrettyHandler(os.Stdout, opts)
		slogInstance = slog.New(handler)
	case "prod":
		// TODO: create log json file
		opts := slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		handler := slog.NewJSONHandler(os.Stdout, &opts)
		slogInstance = slog.New(handler)
	}

	slog.SetDefault(slogInstance)
	return nil, slogInstance
}
