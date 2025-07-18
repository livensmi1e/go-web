package main

import (
	_ "go-web/docs"
	"go-web/internal/platform"
	"go-web/internal/transport/http"
	"log/slog"
)

//	@title			Go Web Service API Document
//	@version		1.0
//	@description	Web Service API Template using Go net/http

// @host		localhost:8000
// @BasePath	/api
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
