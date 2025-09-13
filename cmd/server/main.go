package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"

	_ "go-web/docs"
	"go-web/internal/platform"
	httpTransport "go-web/internal/transport/http"
)

// @title			Go Web Service API Document
// @version		1.0
// @description	Web Service API Template using Go net/http
// @host		localhost:8000
// @BasePath	/api
func main() {
	cfg := platform.NewConfig()
	logger := platform.NewLogger(cfg)
	slog.SetDefault(logger)

	if cfg.MonitorEnabled {
		go func() {
			slog.Info("monitor server running...", "addr", cfg.MonitorServerAddr())
			if err := platform.RunMonitor(cfg.MonitorServerAddr()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("monitor server failed running: ", "error", err.Error())
			}
		}()
	}

	go func() {
		slog.Info("http server running...", "addr", cfg.HttpServerAddr())
		if err := httpTransport.RunServer(cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server failed running: ", "error", err.Error())
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	slog.Info("received interrupt signal, stopping servers...")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := httpTransport.StopServer(); err != nil {
			slog.Error("failed to stop http server: ", "error", err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		if err := platform.StopMonitor(); err != nil {
			slog.Error("failed to stop monitor server: ", "error", err.Error())
		}
	}()

	wg.Wait()
}
