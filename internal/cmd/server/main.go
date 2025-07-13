package main

import (
	"go-web/internal/platform"
	"go-web/internal/transport/http"
	"log/slog"
)

func main() {
	cfg := platform.NewConfig()
	logger := platform.NewLogger(cfg)
	slog.SetDefault(logger)

	slog.Info("server running...", "addr", cfg.HttpServerAddr())
	if err := http.Run(cfg); err != nil {
		slog.Error("server failed running: ", "error", err.Error())
	}
}
