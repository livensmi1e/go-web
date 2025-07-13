package platform

import (
	"log/slog"
	"os"
)

func NewLogger(cfg *Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	if cfg.Env == "prod" {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
