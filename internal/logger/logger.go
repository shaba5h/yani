package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/shaba5h/yani/internal/config"
)

func New(env config.Env) *slog.Logger {
	switch env {
	case config.EnvLocal:
		return slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.Kitchen}))
	case config.EnvDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
