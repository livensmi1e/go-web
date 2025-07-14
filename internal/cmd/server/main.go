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

	if cfg.MonitorEnable {
		go func() {
			slog.Info("monitor server running...", "addr", cfg.MonitorServerAddr())
			if err := platform.RunMonitor(cfg.MonitorServerAddr()); err != nil {
				slog.Error("monitor server failed running: ", "error", err.Error())
			}
		}()
	}

	slog.Info("http server running...", "addr", cfg.HttpServerAddr())
	if err := http.RunServer(cfg); err != nil {
		slog.Error("http server failed running: ", "error", err.Error())
	}
}
